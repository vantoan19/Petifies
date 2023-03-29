import 'dart:async';

import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/comment/repository/comment_repository.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/comment.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/uploading_comment.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'list_comments_controller.g.dart';

@Riverpod(keepAlive: false)
class ListCommentsRefreshNotifier extends _$ListCommentsRefreshNotifier {
  @override
  int build({required String parentID}) {
    return 0;
  }
}

@Riverpod(keepAlive: false)
class CommentStreamRequestController extends _$CommentStreamRequestController {
  @override
  StreamController<UserListCommentsByParentIDRequest>? build(
      {required String parentID}) {
    return null;
  }

  void update(
      StreamController<UserListCommentsByParentIDRequest> requestController) {
    state = requestController;
  }
}

@Riverpod(keepAlive: false)
class ListCommentsController extends _$ListCommentsController {
  @override
  Stream<List<CommentModel>> build({required String parentID}) async* {
    ref.watch(listCommentsRefreshNotifierProvider(parentID: parentID));
    final commentRepository = ref.read(commentRepositoryProvider);

    final commentStream =
        await commentRepository.userListCommentsByParentID(parentID: parentID);
    ref
        .read(
            commentStreamRequestControllerProvider(parentID: parentID).notifier)
        .update(commentStream.first);

    await for (final response in commentStream.second) {
      List<CommentModel> comments = [];
      for (final commentResp in response.comments) {
        final comment = CommentModel(
          id: commentResp.id,
          owner: BasicUserInfoModel(
              id: commentResp.author.id,
              email: commentResp.author.email,
              firstName: commentResp.author.firstName,
              lastName: commentResp.author.lastName),
          createdAt: commentResp.createdAt.toDateTime(),
          postID: commentResp.postId,
          textContent: commentResp.content == "" ? null : commentResp.content,
          image: commentResp.image.uri == ""
              ? null
              : NetworkImageModel(uri: commentResp.image.uri),
          video: commentResp.video.uri == ""
              ? null
              : NetworkVideoModel(uri: commentResp.video.uri),
          hasReacted: commentResp.hasReacted,
          loveCount: commentResp.loveCount,
          subcommentCount: commentResp.subcommentCount,
        );
        comments.add(comment);
      }
      ref
          .read(commentsProvider(parentID: parentID).notifier)
          .addCommentFeeds(comments);
      yield comments;
    }
  }
}

@Riverpod(keepAlive: false)
class NewlyCreatedComments extends _$NewlyCreatedComments {
  @override
  List<UploadingCommentModel> build({required String parentID}) {
    return [];
  }

  void addNewlyCreatedComment(UploadingCommentModel comment) {
    state = [...state, comment];
  }

  void removePost(String tempId) async {
    final idx = state.indexWhere((element) => element.tempId == tempId);

    state = [...state.sublist(0, idx), ...state.sublist(idx + 1)];
  }
}

@Riverpod(keepAlive: false)
class Comments extends _$Comments {
  @override
  List<CommentModel> build({required String parentID}) {
    return [];
  }

  void addCommentToHead(CommentModel comment) {
    state = [comment, ...state];
  }

  void addCommentFeed(CommentModel comment) {
    state = [...state, comment];
  }

  void addCommentFeeds(List<CommentModel> comments) {
    state = [...state, ...comments];
  }

  void reset() {
    state = [];
  }
}
