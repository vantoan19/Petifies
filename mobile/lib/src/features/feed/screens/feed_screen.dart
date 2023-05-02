// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/feed/controller/newfeed_controller.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/uploading_post.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/widgets/appbars/main_appbar.dart';
import 'package:mobile/src/widgets/floating_buttons/new_post_floating_button.dart';
import 'package:mobile/src/widgets/posts/post.dart';
import 'package:mobile/src/widgets/posts/uploading_post.dart';
import 'package:mobile/src/widgets/story_circle/story_circle.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'feed_screen.g.dart';

@Riverpod(keepAlive: true)
class NewlyCreatedPosts extends _$NewlyCreatedPosts {
  @override
  List<UploadingPostModel> build() {
    return [];
  }

  void addNewlyCreatedPost(UploadingPostModel post) {
    state = [...state, post];
  }

  void removePost(String tempId) async {
    final idx = state.indexWhere((element) => element.tempId == tempId);

    state = [...state.sublist(0, idx), ...state.sublist(idx + 1)];
  }
}

@Riverpod(keepAlive: true)
class PostFeeds extends _$PostFeeds {
  @override
  List<PostModel> build() {
    return [];
  }

  void addPostToHead(PostModel post) {
    state = [post, ...state];
  }

  void addPostFeed(PostModel post) {
    state = [...state, post];
  }

  void addPostFeeds(List<PostModel> posts) {
    state = [...state, ...posts];
  }

  void reset() {
    state = [];
  }
}

class FeedScreen extends ConsumerWidget {
  const FeedScreen();

  // Widget get _stories {
  //   return SizedBox(
  //     height: 100,
  //     child: ListView(
  //       scrollDirection: Axis.horizontal,
  //       children: [
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //         StoryCircle(),
  //       ],
  //     ),
  //   );
  // }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);
    final newfeed = ref.watch(newFeedStreamControllerProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    final isFeedLoading = newfeed.isLoading;
    final newlyCreatedPost = ref.watch(newlyCreatedPostsProvider);
    final postFeeds = ref.watch(postFeedsProvider);

    return Scaffold(
      appBar: MainAppBar(),
      floatingActionButton: NewPostFloatingButton(
        heroTag: "feedTag",
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.separated(
              separatorBuilder: (context, index) => Divider(
                thickness: 5,
                color: Theme.of(context).colorScheme.tertiary,
              ),
              itemBuilder: (context, index) {
                if (index == 0) {
                  return SizedBox.shrink();
                }
                if (newlyCreatedPost.length > 0 &&
                    index <= newlyCreatedPost.length) {
                  return ProviderScope(
                    key: ObjectKey(newlyCreatedPost[index - 1].tempId),
                    overrides: [
                      uploadingPostInfoProvider
                          .overrideWithValue(newlyCreatedPost[index - 1]),
                    ],
                    child: UploadingPost(),
                  );
                }
                return ProviderScope(
                  key: ObjectKey(
                      postFeeds[index - newlyCreatedPost.length - 1].id),
                  overrides: [
                    postInfoProvider.overrideWithValue(
                        postFeeds[index - newlyCreatedPost.length - 1]),
                    isPostContextProvider.overrideWithValue(true),
                  ],
                  child: Post(),
                );
              },
              itemCount: newlyCreatedPost.length + postFeeds.length + 1,
            ),
          ),
          if (isFeedLoading) CircularProgressIndicator(),
        ],
      ),
    );
  }
}
