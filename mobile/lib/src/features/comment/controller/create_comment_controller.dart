import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/comment/controller/comment_count_controller.dart';
import 'package:mobile/src/features/comment/controller/list_comments_controller.dart';
import 'package:mobile/src/features/comment/repository/comment_repository.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/comment.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/uploading_comment.dart';
import 'package:mobile/src/models/video.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'create_comment_controller.g.dart';

@Riverpod(keepAlive: false)
class CreateCommentController extends _$CreateCommentController {
  @override
  void build() {
    return;
  }

  Future<Either<Failure, CommentModel?>> createComment({
    required String tempId,
    required BasicUserInfoModel author,
    required String postID,
    required String parentID,
    required bool isPostParent,
    String textContent = "",
    Future<Either<Failure, NetworkImageModel>>? imageFuture = null,
    Future<Either<Failure, NetworkVideoModel>>? videoFuture = null,
  }) async {
    final commentRepository = ref.read(commentRepositoryProvider);

    NetworkImageModel? image = null;
    NetworkVideoModel? video = null;

    if (imageFuture != null) {
      final imageEither = await imageFuture;
      imageEither.fold((l) => null, (r) => image = r);
    }

    if (videoFuture != null) {
      final videoEither = await videoFuture;
      videoEither.fold((l) => null, (r) => video = r);
    }

    final commentFuture = await commentRepository.createComment(
      author: author,
      postID: postID,
      parentID: parentID,
      isPostParent: isPostParent,
      textContent: textContent,
      image: image,
      video: video,
    );

    commentFuture.fold((l) => null, (r) {
      ref
          .read(newlyCreatedCommentsProvider(parentID: parentID).notifier)
          .removePost(tempId);
      ref
          .read(commentsProvider(parentID: parentID).notifier)
          .addCommentToHead(r);
      ref
          .read(commentCountProvider(Tuple2(parentID, isPostParent)).notifier)
          .incr();
    });

    return commentFuture;
  }
}
