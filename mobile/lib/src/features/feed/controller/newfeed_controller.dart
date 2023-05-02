import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/feed/repository/newfeed_repository.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/auth-gateway/v1/auth-gateway.v1.pbgrpc.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'newfeed_controller.g.dart';

@riverpod
class NewFeedRefreshNotifier extends _$NewFeedRefreshNotifier {
  @override
  int build() {
    return 0;
  }

  void refreshFeed() {
    ref.read(postFeedsProvider.notifier).reset();
    state = state + 1;
  }
}

@riverpod
class NewFeedStreamRequestController extends _$NewFeedStreamRequestController {
  @override
  StreamController<ListNewFeedsRequest>? build() {
    return null;
  }

  void update(StreamController<ListNewFeedsRequest> controller) {
    state = controller;
  }
}

@riverpod
class NewFeedStreamController extends _$NewFeedStreamController {
  @override
  Stream<List<PostModel>> build() async* {
    ref.watch(newFeedRefreshNotifierProvider);
    final newfeedRepository = ref.read(newfeedRepositoryProvider);

    final newfeedStream = await newfeedRepository.getListNewfeedStream();
    ref
        .read(newFeedStreamRequestControllerProvider.notifier)
        .update(newfeedStream.first);

    await for (final response in newfeedStream.second) {
      List<PostModel> posts = [];
      for (final postResp in response.posts) {
        final post = PostModel(
          id: postResp.id,
          owner: BasicUserInfoModel(
              id: postResp.author.id,
              email: postResp.author.email,
              firstName: postResp.author.firstName,
              lastName: postResp.author.lastName),
          postActivity: postResp.activity,
          createdAt: postResp.createdAt.toDateTime(),
          textContent: postResp.content == "" ? null : postResp.content,
          images: postResp.images
              .map((e) => NetworkImageModel(uri: e.uri))
              .toList(),
          videos: postResp.videos
              .map((e) => NetworkVideoModel(uri: e.uri))
              .toList(),
          hasReacted: postResp.hasReacted,
          loveCount: postResp.loveCount,
          commentCount: postResp.commentCount,
        );
        posts.add(post);
      }
      ref.read(postFeedsProvider.notifier).addPostFeeds(posts);
      yield posts;
    }
  }

  void fetchMoreNewFeed() {
    ref
        .read(newFeedStreamRequestControllerProvider)
        ?.add(ListNewFeedsRequest());
  }
}
