import 'package:fpdart/fpdart.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/features/comment/repository/comment_repository.dart';
import 'package:mobile/src/features/post/repository/post_repository.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'comment_count_controller.g.dart';

@Riverpod(keepAlive: false)
class HasChangedCommentCount extends _$HasChangedCommentCount {
  @override
  bool build(Tuple2<String, bool> target) {
    return false;
  }
}

@Riverpod(keepAlive: false)
class CommentCount extends _$CommentCount {
  @override
  int build(Tuple2<String, bool> target) {
    return 0;
  }

  void setCommentCount(int count) {
    state = count;
    setHasChangedCommentCount();
  }

  void desc() {
    state = state - 1;
    setHasChangedCommentCount();
  }

  void incr() {
    state = state + 1;
    setHasChangedCommentCount();
  }

  void setHasChangedCommentCount() {
    if (ref.read(hasChangedCommentCountProvider(this.target).notifier).state ==
        false) {
      ref.read(hasChangedCommentCountProvider(this.target).notifier).state =
          true;
    }
  }
}

@Riverpod(keepAlive: false)
class CommentCountController extends _$CommentCountController {
  @override
  Stream<int> build(Tuple2<String, bool> target) async* {
    ResponseStream<StreamCommentCountResponse> commentCountStream;
    if (target.second) {
      final postRepository = ref.read(postRepositoryProvider);
      commentCountStream =
          await postRepository.getPostCommentCountStream(postID: target.first);
    } else {
      final commentRepository = ref.read(commentRepositoryProvider);
      commentCountStream = await commentRepository
          .getCommentSubCommentCountStream(commentID: target.first);
    }

    await for (final response in commentCountStream) {
      ref
          .read(commentCountProvider(target).notifier)
          .setCommentCount(response.commentCount.value);
      yield response.commentCount.value;
    }
  }
}
