import 'dart:async';

import 'package:fpdart/fpdart.dart';
import 'package:grpc/grpc.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/comment/repository/comment_repository.dart';
import 'package:mobile/src/features/post/repository/post_repository.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pb.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'love_count_controller.g.dart';

@Riverpod(keepAlive: false)
class HasChangedLoveCount extends _$HasChangedLoveCount {
  @override
  bool build(Tuple2<String, bool> target) {
    return false;
  }
}

@Riverpod(keepAlive: false)
class LoveCount extends _$LoveCount {
  @override
  int build(Tuple2<String, bool> target) {
    return 0;
  }

  void setLoveCount(int count) {
    state = count;
    setHasChangedLoveCount();
  }

  void desc() {
    state = state - 1;
    setHasChangedLoveCount();
  }

  void incr() {
    state = state + 1;
    setHasChangedLoveCount();
  }

  void setHasChangedLoveCount() {
    if (ref.read(hasChangedLoveCountProvider(this.target).notifier).state ==
        false) {
      ref.read(hasChangedLoveCountProvider(this.target).notifier).state = true;
    }
  }
}

@Riverpod(keepAlive: false)
class HasReacted extends _$HasReacted {
  @override
  bool? build(Tuple2<String, bool> target) {
    return null;
  }

  void setHasReacted(bool hasReacted) {
    state = hasReacted;
  }
}

@Riverpod(keepAlive: false)
class LoveCountController extends _$LoveCountController {
  @override
  Stream<int> build(Tuple2<String, bool> target) async* {
    ResponseStream<StreamLoveCountResponse> loveCountStream;
    if (target.second) {
      final postRepository = ref.read(postRepositoryProvider);
      loveCountStream =
          await postRepository.getPostLoveCountStream(postID: target.first);
    } else {
      final commentRepository = ref.read(commentRepositoryProvider);
      loveCountStream = await commentRepository.getCommentLoveCountStream(
          commentID: target.first);
    }

    await for (final response in loveCountStream) {
      ref
          .read(loveCountProvider(target).notifier)
          .setLoveCount(response.loveCount.value);
      yield response.loveCount.value;
    }
  }

  Future<void> toggleLoveReact() async {
    Either<Failure, bool> result;
    if (this.target.second) {
      final postRepository = ref.read(postRepositoryProvider);
      result = await postRepository.userToggleLoveReactPost(
          postID: this.target.first);
    } else {
      final commentRepository = ref.read(commentRepositoryProvider);
      result = await commentRepository.userToggleLoveReactComment(
          commentID: this.target.first);
    }
    result.fold((l) => throw l, (r) {
      ref.read(hasReactedProvider(this.target).notifier).setHasReacted(r);
      if (r) {
        ref.read(loveCountProvider(this.target).notifier).incr();
      } else {
        ref.read(loveCountProvider(this.target).notifier).desc();
      }
    });
  }
}
