import 'dart:async';
import 'dart:ffi';
import 'dart:io';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/post/repository/file_repository.dart';
import 'package:mobile/src/features/post/repository/post_repository.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/models/video.dart';

final createPostControllerProvider =
    AsyncNotifierProvider.autoDispose<CreatePostController, void>(
        CreatePostController.new);

class CreatePostController extends AutoDisposeAsyncNotifier<void> {
  @override
  void build() {
    return;
  }

  Future<Either<Failure, PostModel?>> createPost({
    required UserModel author,
    required String textContent,
    required List<NetworkImageModel> images,
    required List<NetworkVideoModel> videos,
  }) async {
    final postRepository = ref.read(postRepositoryProvider);

    final postFuture = postRepository.createPost(
      author: author,
      textContent: textContent,
      images: images,
      videos: videos,
    );

    ref
        .read(newlyCreatedPostFuturesProvider.notifier)
        .addPostFuture(postFuture);
    ref.read(newlyCreatedPostProvider.notifier).addNewlyCreatedPost(
          PostModel(
            owner: author,
            postActivity: "post",
            textContent: textContent,
            images: images,
            videos: videos,
            createdAt: DateTime.now(),
            loveCount: 0,
            commentCount: 0,
          ),
        );

    return postFuture;
  }
}
