import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/google/protobuf/timestamp.pb.dart';
import 'package:mobile/src/providers/model_providers.dart';
import 'package:mobile/src/widgets/posts/post.dart';
import 'package:mobile/src/widgets/story_circle/story_circle.dart';

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

    return ListView(
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
        Post(
          postData: PostModel(
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
            postTime: Timestamp(),
            textContent: "Miaomao",
            images: [
              NetworkImageModel(
                uri: "https://picsum.photos/200/300",
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
        )
      ],
    );
  }
}
