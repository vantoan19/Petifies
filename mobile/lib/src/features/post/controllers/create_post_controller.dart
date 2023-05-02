import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/post/repository/post_repository.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/video.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'create_post_controller.g.dart';

@Riverpod(keepAlive: false)
class CreatePostController extends _$CreatePostController {
  @override
  void build() {
    return;
  }

  Future<Either<Failure, PostModel?>> createPost({
    required String tempId,
    required BasicUserInfoModel author,
    required String textContent,
    required String activity,
    String visibility = "public",
    required List<Future<Either<Failure, NetworkImageModel>>> imageFutures,
    required List<Future<Either<Failure, NetworkVideoModel>>> videoFutures,
  }) async {
    final postRepository = ref.read(postRepositoryProvider);

    // wait for uploading images and videos to be fullfilled
    final imageEithers = await Future.wait(imageFutures);
    final videoEithers = await Future.wait(videoFutures);
    List<NetworkImageModel> images = [];
    List<NetworkVideoModel> videos = [];
    imageEithers.forEach((either) {
      either.fold((l) => null, (r) => images.add(r));
    });
    videoEithers.forEach((either) {
      either.fold((l) => null, (r) => videos.add(r));
    });

    final postEither = await postRepository.createPost(
      author: author,
      visibility: visibility,
      activity: activity,
      textContent: textContent,
      images: images,
      videos: videos,
    );

    postEither.fold((l) => null, (r) {
      ref.read(newlyCreatedPostsProvider.notifier).removePost(tempId);
      ref.read(postFeedsProvider.notifier).addPostToHead(r);
    });

    return postEither;
  }
}
