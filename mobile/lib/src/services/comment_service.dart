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
import 'package:flutter_riverpod/flutter_riverpod.dart';

class CommentService {
  AuthGatewayClient? _authClientInstance;
  Ref _ref;

  CommentService({
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

  Future<CommentWithUserInfo> userCreateComment({
    required String postID,
    required String parentID,
    required bool isPostParent,
    required String textContent,
    NetworkImageModel? image = null,
    NetworkVideoModel? video = null,
  }) async {
    UserCreateCommentRequest request = UserCreateCommentRequest(
      postId: postID,
      parentId: parentID,
      isPostParent: isPostParent,
      content: textContent,
      image: (image != null
          ? Image(uri: image.uri, description: image.description)
          : Image()),
      video: (video != null
          ? Video(uri: video.uri, description: video.description)
          : Video()),
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userCreateComment(request), _ref);
  }

  Future<UserToggleLoveResponse> userToggleLoveReactComment({
    required String commentID,
  }) async {
    UserToggleLoveRequest request = UserToggleLoveRequest(
      targetId: commentID,
      isPostTarget: false,
    );

    return await RetryUtils.RetryRefreshToken(
        () async => (await _authClient).userToggleLoveReact(request), _ref);
  }

  Future<
          Tuple2<StreamController<UserListCommentsByParentIDRequest>,
              ResponseStream<UserListCommentsByParentIDResponse>>>
      userListCommentsByParentID({
    required String parentID,
    required int pageSize,
  }) async {
    final requestController =
        StreamController<UserListCommentsByParentIDRequest>();
    requestController.add(UserListCommentsByParentIDRequest(
        parentId: parentID, pageSize: pageSize));
    final responseStream = await RetryUtils.RetryRefreshToken(
      () async => (await _authClient)
          .userListCommentsByParentID(requestController.stream),
      _ref,
    );

    return Tuple2(requestController, responseStream);
  }

  Future<ResponseStream<StreamLoveCountResponse>> getLoveCount({
    required String commentID,
  }) async {
    StreamLoveCountRequest request = StreamLoveCountRequest(
      targetId: commentID,
      isPostTarget: false,
    );

    return await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).streamLoveCount(request),
      _ref,
    );
  }

  Future<ResponseStream<StreamCommentCountResponse>> getSubCommentCount({
    required String commentID,
  }) async {
    StreamCommentCountRequest request = StreamCommentCountRequest(
      parentId: commentID,
      isPostParent: false,
    );

    return await RetryUtils.RetryRefreshToken(
      () async => (await _authClient).streamCommentCount(request),
      _ref,
    );
  }
}
