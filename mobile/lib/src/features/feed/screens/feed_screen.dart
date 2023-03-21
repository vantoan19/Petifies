import 'dart:ffi';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/posts/post.dart';
import 'package:mobile/src/widgets/story_circle/story_circle.dart';

final newlyCreatedPostFuturesProvider = NotifierProvider<
    NewlyCreatedPostFutures,
    List<Future<Either<Failure, PostModel?>>>>(NewlyCreatedPostFutures.new);

final newlyCreatedPostProvider =
    NotifierProvider<NewlyCreatedPosts, List<PostModel>>(NewlyCreatedPosts.new);

final postFeedsProvider =
    NotifierProvider<PostFeeds, List<PostModel>>(PostFeeds.new);

class NewlyCreatedPostFutures
    extends Notifier<List<Future<Either<Failure, PostModel?>>>> {
  @override
  List<Future<Either<Failure, PostModel?>>> build() {
    return [];
  }

  void addPostFuture(Future<Either<Failure, PostModel?>> postFuture) {
    state = [postFuture, ...state];
  }

  void removePostFuture(int idx) {
    state = [...state.sublist(0, idx), ...state.sublist(idx + 1)];
  }
}

class NewlyCreatedPosts extends Notifier<List<PostModel>> {
  @override
  List<PostModel> build() {
    return [];
  }

  void addNewlyCreatedPost(PostModel post) {
    state = [...state, post];
  }

  void removePost(int idx) async {
    final postData = await ref.read(newlyCreatedPostFuturesProvider)[idx];
    state = [...state.sublist(0, idx), ...state.sublist(idx + 1)];
    ref.read(newlyCreatedPostFuturesProvider.notifier).removePostFuture(idx);
    postData.fold((l) => null, (r) {
      ref.read(postFeedsProvider.notifier).addPostToHead(r!);
    });
  }
}

class PostFeeds extends Notifier<List<PostModel>> {
  @override
  List<PostModel> build() {
    return [
      PostModel(
        owner: UserModel(
          email: "toan@gmail.com",
          id: "123",
          firstName: "Toan",
          lastName: "Tran",
          isActivated: true,
          countPost: 0,
          followers: 0,
          following: 0,
        ),
        postActivity: "has just shared a new post",
        createdAt: DateTime.now().subtract(Duration(days: 31)),
        textContent: "Miaomao",
        images: [
          NetworkImageModel(
            uri:
                "https://storage.googleapis.com/petifies-storage/1717e73c-5811-48d6-b532-8fb173cc45d9-image_picker_772C4BDD-A3F7-45E3-AD31-81A81A1FC59D-67222-000012CF3A3B625F.jpg",
          )
        ],
        videos: [
          NetworkVideoModel(
            uri:
                "https://storage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4",
          )
        ],
        loveCount: 100000,
        commentCount: 10,
      ),
    ];
  }

  void addPostToHead(PostModel post) {
    state = [post, ...state];
  }

  void addPostFeed(PostModel post) {
    state = [...state, post];
  }
}

class FeedScreen extends ConsumerWidget {
  const FeedScreen({super.key});

  Widget get _stories {
    return SizedBox(
      height: 100,
      child: ListView(
        scrollDirection: Axis.horizontal,
        children: [
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
          StoryCircle(),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    final newlyCreatedPost = ref.watch(newlyCreatedPostProvider);
    final newlyCreatedPostFuture = ref.watch(newlyCreatedPostFuturesProvider);
    final postFeeds = ref.watch(postFeedsProvider);

    return Column(
      children: [
        Container(
          height: 5,
          color: Color.fromRGBO(10, 10, 10, 10),
        ),
        _stories,
        Container(
          height: 5,
          color: Color.fromRGBO(10, 10, 10, 10),
        ),
        Expanded(
          child: ListView.separated(
            separatorBuilder: (context, index) => Divider(
              thickness: 5,
              color: Theme.of(context).colorScheme.tertiary,
            ),
            itemBuilder: (context, index) {
              if (newlyCreatedPost.length > 0 &&
                  index < newlyCreatedPost.length) {
                return Column(
                  children: [
                    FutureBuilder(
                      future: newlyCreatedPostFuture[index],
                      builder: (context, snapshot) {
                        if (snapshot.hasData) {
                          ref
                              .read(newlyCreatedPostProvider.notifier)
                              .removePost(index);
                          return Text("Post created successfully");
                        } else if (snapshot.hasError) {
                          return Text("Failed to create post");
                        } else {
                          return LinearProgressIndicator();
                        }
                      },
                    ),
                    Container(
                      foregroundDecoration: BoxDecoration(
                        color: Themes.darkGreyColorTransparent,
                      ),
                      child: Post(postData: newlyCreatedPost[index]),
                    ),
                  ],
                );
                // } else if (index == newlyCreatedPost.length + postFeeds.length) {
                //   return Padding(
                //     padding: const EdgeInsets.fromLTRB(0, 0, 0, 30),
                //     child: Center(
                //       child: Column(
                //         children: [
                //           Image.asset(
                //             Constants.emptyBoxPng,
                //             width: MediaQuery.of(context).size.width * 0.3,
                //           ),
                //           Text("No more feeds to load"),
                //         ],
                //       ),
                //     ),
                //   );
              }
              return Post(postData: postFeeds[index - newlyCreatedPost.length]);
            },
            itemCount: newlyCreatedPost.length + postFeeds.length,
          ),
        )
      ],
    );
  }
}
