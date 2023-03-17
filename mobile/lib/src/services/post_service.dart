import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/services/grpc/auth_interceptor.dart';
import 'package:mobile/src/services/grpc/grpc_flutter_client.dart';
import 'package:mobile/src/utils/retryutils.dart';
import 'package:riverpod/riverpod.dart';

class PostService {
  AuthGatewayClient? _authClientInstance;
  Ref _ref;

  PostService({
    required ref,
  }) : this._ref = ref;

  Future<AuthGatewayClient> get _authClient async {
    if (_authClientInstance == null) {
      _authClientInstance = AuthGatewayClient(
          await GrpcFlutterClient.getClient(),
          interceptors: [AuthInterceptor(ref: _ref)]);
    }
    return _authClientInstance!;
  }

  Future<Post> userCreatePost({
    required String textContent,
    required List<NetworkImageModel> images,
    required List<NetworkVideoModel> videos,
  }) async {
    UserCreatePostRequest request = UserCreatePostRequest(
      content: textContent,
      images: images.map((image) => Image(uri: image.uri, description: "")),
      videos: videos.map((video) => Video(uri: video.uri, description: "")),
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userCreatePost(request), _ref);
  }
}
