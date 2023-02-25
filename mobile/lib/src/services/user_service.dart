import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/proto/public-gateway/v1/public-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/model_providers.dart';
import 'package:mobile/src/providers/storages_provider.dart';
import 'package:mobile/src/services/grpc/auth_interceptor.dart';
import 'package:mobile/src/services/grpc/grpc_flutter_client.dart';
import 'package:retry/retry.dart';

class UserService {
  PublicGatewayClient? _publicClientInstance;
  AuthGatewayClient? _authClientInstance;
  Ref _ref;

  UserService({required ref}) : _ref = ref;

  Future<PublicGatewayClient> get _publicClient async {
    if (_publicClientInstance == null) {
      _publicClientInstance =
          PublicGatewayClient(await GrpcFlutterClient.getClient());
    }
    return _publicClientInstance!;
  }

  Future<AuthGatewayClient> get _authClient async {
    if (_authClientInstance == null) {
      _authClientInstance = AuthGatewayClient(
          await GrpcFlutterClient.getClient(),
          interceptors: [AuthInterceptor(ref: _ref)]);
    }
    return _authClientInstance!;
  }

  Future<User> createUser({
    required String email,
    required String password,
    required String firstName,
    required String lastName,
  }) async {
    CreateUserRequest request = CreateUserRequest(
      email: email,
      password: password,
      firstName: firstName,
      lastName: lastName,
    );

    return (await _publicClient).createUser(request);
  }

  Future<LoginResponse> login({
    required String email,
    required String password,
  }) async {
    LoginRequest request = LoginRequest(
      email: email,
      password: password,
    );

    return (await _publicClient).login(request);
  }

  Future<RefreshTokenResponse> refreshToken({
    required String refreshToken,
  }) async {
    RefreshTokenRequest request =
        RefreshTokenRequest(refreshToken: refreshToken);

    return (await _publicClient).refreshToken(request);
  }

  Future<User> getMyInfo() async {
    GetMyInfoRequest request = GetMyInfoRequest();

    return await _retryRefreshToken(
      () async => (await _authClient).getMyInfo(request),
    );
  }

  Future<dynamic> _retryRefreshToken(Function fn) async {
    return retry(
      () => fn(),
      maxAttempts: 2,
      delayFactor: Duration.zero,
      retryIf: (e) =>
          e is GrpcError &&
          e.message == "token has expired" &&
          e.code == StatusCode.unauthenticated,
      onRetry: (e) async {
        final secureStorage = _ref.read(secureStorageProvider);
        final String? rfToken = await secureStorage.read(key: "refreshToken");
        debugPrint(rfToken);
        try {
          final accessToken =
              (await refreshToken(refreshToken: rfToken!)).accessToken;
          await secureStorage.write(key: "accessToken", value: accessToken);
        } catch (e) {
          _ref.read(myUserProvider.notifier).SetUser(null);
        }
      },
    );
  }
}
