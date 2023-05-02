import 'dart:async';

import 'package:fpdart/fpdart.dart';
import 'package:grpc/service_api.dart';
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
        interceptors: [AuthInterceptor(ref: _ref)],
      );
    }
    return _authClientInstance!;
  }

  Future<PostWithUserInfo> userCreatePost({
    required String visibility,
    required String activity,
    required String textContent,
    required List<NetworkImageModel> images,
    required List<NetworkVideoModel> videos,
  }) async {
    UserCreatePostRequest request = UserCreatePostRequest(
      visibility: visibility,
      activity: activity,
      content: textContent,
      images: images.map((image) => Image(uri: image.uri, description: "")),
      videos: videos.map((video) => Video(uri: video.uri, description: "")),
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userCreatePost(request), _ref);
  }

  Future<UserToggleLoveResponse> userToggleLoveReactPost({
    required String postID,
  }) async {
    UserToggleLoveRequest request = UserToggleLoveRequest(
      targetId: postID,
      isPostTarget: true,
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userToggleLoveReact(request), _ref);
  }

  Future<
      Tuple2<StreamController<ListNewFeedsRequest>,
          ResponseStream<ListNewFeedsResponse>>> listNewFeeds() async {
    final requestController = StreamController<ListNewFeedsRequest>();
    requestController.add(ListNewFeedsRequest());
    final responseStream = await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).listNewFeeds(requestController.stream),
      _ref,
    );

    return Tuple2(requestController, responseStream);
  }

  Future<ResponseStream<StreamLoveCountResponse>> getLoveCount({
    required String postID,
  }) async {
    StreamLoveCountRequest request = StreamLoveCountRequest(
      targetId: postID,
      isPostTarget: true,
    );

    return await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).streamLoveCount(request),
      _ref,
    );
  }

  Future<ResponseStream<StreamCommentCountResponse>> getCommentCount({
    required String postID,
  }) async {
    StreamCommentCountRequest request = StreamCommentCountRequest(
      parentId: postID,
      isPostParent: true,
    );

    return await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).streamCommentCount(request),
      _ref,
    );
  }
}
