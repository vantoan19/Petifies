import 'dart:async';

import 'package:fpdart/fpdart.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/comment.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/comment_service.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

final commentRepositoryProvider = Provider<ICommentRepository>((ref) =>
    CommentRepository(commentService: ref.read(commentServiceProvider)));

abstract class ICommentRepository {
  Future<Either<Failure, CommentModel>> createComment({
    required BasicUserInfoModel author,
    required String postID,
    required String parentID,
    required bool isPostParent,
    required String textContent,
    NetworkImageModel? image = null,
    NetworkVideoModel? video = null,
  });
  Future<ResponseStream<StreamLoveCountResponse>> getCommentLoveCountStream({
    required String commentID,
  });
  Future<ResponseStream<StreamCommentCountResponse>>
      getCommentSubCommentCountStream({
    required String commentID,
  });
  Future<Either<Failure, bool>> userToggleLoveReactComment({
    required String commentID,
  });
  Future<
          Tuple2<StreamController<UserListCommentsByParentIDRequest>,
              ResponseStream<UserListCommentsByParentIDResponse>>>
      userListCommentsByParentID({
    required String parentID,
  });
}

class CommentRepository implements ICommentRepository {
  final CommentService _commentService;

  CommentRepository({required CommentService commentService})
      : this._commentService = commentService;

  Future<Either<Failure, CommentModel>> createComment({
    required BasicUserInfoModel author,
    required String postID,
    required String parentID,
    required bool isPostParent,
    required String textContent,
    NetworkImageModel? image = null,
    NetworkVideoModel? video = null,
  }) async {
    try {
      CommentWithUserInfo resp = await _commentService.userCreateComment(
        postID: postID,
        parentID: parentID,
        isPostParent: isPostParent,
        textContent: textContent,
        image: image,
      );

      CommentModel comment = CommentModel(
        id: resp.id,
        owner: author,
        postID: resp.postId,
        parentID: resp.parentId,
        isPostParent: resp.isPostParent,
        createdAt: resp.createdAt.toDateTime(),
        textContent: resp.content == "" ? null : resp.content,
        image: resp.image.uri == ""
            ? null
            : NetworkImageModel(
                uri: resp.image.uri,
                description: resp.image.description,
              ),
        video: resp.video.uri == ""
            ? null
            : NetworkVideoModel(
                uri: resp.video.uri,
                description: resp.video.description,
              ),
        hasReacted: resp.hasReacted,
        loveCount: 0,
        subcommentCount: 0,
      );
      return right(comment);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<ResponseStream<StreamLoveCountResponse>> getCommentLoveCountStream({
    required String commentID,
  }) {
    return _commentService.getLoveCount(commentID: commentID);
  }

  Future<ResponseStream<StreamCommentCountResponse>>
      getCommentSubCommentCountStream({
    required String commentID,
  }) {
    return _commentService.getSubCommentCount(commentID: commentID);
  }

  Future<Either<Failure, bool>> userToggleLoveReactComment({
    required String commentID,
  }) async {
    try {
      UserToggleLoveResponse resp = await _commentService
          .userToggleLoveReactComment(commentID: commentID);

      return right(resp.hasReacted.value);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<
          Tuple2<StreamController<UserListCommentsByParentIDRequest>,
              ResponseStream<UserListCommentsByParentIDResponse>>>
      userListCommentsByParentID({
    required String parentID,
  }) async {
    Tuple2<StreamController<UserListCommentsByParentIDRequest>,
            ResponseStream<UserListCommentsByParentIDResponse>> stream =
        await _commentService.userListCommentsByParentID(
            parentID: parentID, pageSize: 40);

    return stream;
  }
}
