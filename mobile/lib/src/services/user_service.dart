import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/proto/public-gateway/v1/public-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/services/grpc_flutter_client.dart';

class UserService {
  PublicGatewayClient? _clientInstance;

  Future<PublicGatewayClient> get _client async {
    if (_clientInstance == null) {
      _clientInstance =
          PublicGatewayClient(await GrpcFlutterClient.getClient());
    }
    return _clientInstance!;
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

    return await (await _client).createUser(request);
  }

  Future<LoginResponse> login({
    required String email,
    required String password,
  }) async {
    LoginRequest request = LoginRequest(
      email: email,
      password: password,
    );

    return await (await _client).login(request);
  }
}
