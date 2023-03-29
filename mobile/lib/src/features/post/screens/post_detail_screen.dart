// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/comment/controller/comment_count_controller.dart';
import 'package:mobile/src/features/comment/controller/list_comments_controller.dart';
import 'package:mobile/src/features/love/controllers/love_count_controller.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/appbars/post_detail_appbar.dart';
import 'package:mobile/src/widgets/buttons/comment_button.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/comment/uploading_comment.dart';
import 'package:mobile/src/widgets/media_view/media_view.dart';
import 'package:mobile/src/widgets/posts/post.dart';
import 'package:mobile/src/widgets/posts/post_head.dart';
import 'package:mobile/src/widgets/textfield/reply_text_field.dart';

class PostDetailScreen extends ConsumerWidget {
  final bool autoFocus;
  const PostDetailScreen({
    required this.autoFocus,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final postData = ref.watch(postInfoProvider);
    final _ = ref.watch(listCommentsControllerProvider(parentID: postData.id));
    final comments = ref.watch(commentsProvider(parentID: postData.id));
    final newlyCreatedComments =
        ref.watch(newlyCreatedCommentsProvider(parentID: postData.id));

    return Scaffold(
      appBar: PostDetailsAppBar(),
      body: Column(
        children: [
          Expanded(
            child: ListView.separated(
              itemBuilder: (context, index) {
                if (index == 0) {
                  return const PostHead(
                    isUploadingPost: false,
                  );
                } else if (index == 1) {
                  return PostDetailsBody(
                    postData: postData,
                  );
                } else {
                  if (newlyCreatedComments.length > 0 &&
                      index - 2 < newlyCreatedComments.length) {
                    return ProviderScope(
                      key: ObjectKey(newlyCreatedComments[index - 2].tempId),
                      overrides: [
                        uploadingCommentInfoProvider
                            .overrideWithValue(newlyCreatedComments[index - 2])
                      ],
                      child: const UploadingComment(),
                    );
                  }
                  return ProviderScope(
                    key: ObjectKey(
                        comments[index - newlyCreatedComments.length - 2].id),
                    overrides: [
                      commentInfoProvider.overrideWithValue(
                          comments[index - newlyCreatedComments.length - 2]),
                      isPostContextProvider.overrideWithValue(false),
                    ],
                    child: const Comment(),
                  );
                }
              },
              separatorBuilder: (context, index) {
                if (index == 0) {
                  return SizedBox.shrink();
                }
                return Divider(
                  thickness: 0.3,
                  color: Theme.of(context).colorScheme.secondary,
                );
              },
              itemCount: comments.length + 2,
            ),
          ),
          Container(
            decoration: BoxDecoration(
              border: Border(
                top: BorderSide(
                  width: 0.25,
                  color: Theme.of(context).colorScheme.secondary,
                ),
              ),
            ),
            padding: EdgeInsets.fromLTRB(0, 8, 0, 12),
            child: ReplyTextField(
              autoFocus: autoFocus,
            ),
          ),
        ],
      ),
    );
  }
}

class PostDetailsBody extends StatelessWidget {
  final PostModel postData;

  const PostDetailsBody({
    Key? key,
    required this.postData,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 24, 0, 8),
      child: Column(
        mainAxisSize: MainAxisSize.max,
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Text Content
              if (postData.textContent != null)
                Padding(
                  padding: const EdgeInsets.fromLTRB(
                      Constants.horizontalScreenPadding + 4,
                      0,
                      Constants.horizontalScreenPadding,
                      24),
                  child: Align(
                    child: Text(
                      postData.textContent!,
                      style: TextStyle(
                        fontSize: 24,
                      ),
                    ),
                    alignment: Alignment.topLeft,
                  ),
                ),
              // Image & Video content
              if ((postData.images != null && postData.images!.length > 0) ||
                  (postData.videos != null && postData.videos!.length > 0))
                MediaView(
                  imageUrls: postData.images != null
                      ? postData.images!.map((e) => e.uri).toList()
                      : [],
                  videoUrls: postData.videos != null
                      ? postData.videos!.map((e) => e.uri).toList()
                      : [],
                  isClickable: true,
                ),
              Column(
                children: [
                  Container(
                    padding: const EdgeInsets.fromLTRB(
                      Constants.horizontalScreenPadding,
                      28,
                      Constants.horizontalScreenPadding,
                      16,
                    ),
                    decoration: BoxDecoration(
                      border: Border(
                        bottom: BorderSide(
                          color: Theme.of(context).colorScheme.secondary,
                          width: 0.3,
                        ),
                      ),
                    ),
                    width: MediaQuery.of(context).size.width * 0.96,
                    alignment: Alignment.centerLeft,
                    child: Container(
                      constraints: BoxConstraints(maxWidth: 240, minWidth: 200),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          _LoveCount(
                            postID: postData.id,
                            initialCount: postData.loveCount,
                          ),
                          _CommentCount(
                            postID: postData.id,
                            initialCount: postData.commentCount,
                          ),
                        ],
                      ),
                    ),
                  ),
                  _ActionButtons(
                    postData: postData,
                  ),
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _LoveCount extends ConsumerWidget {
  final String postID;
  final int initialCount;

  const _LoveCount({
    Key? key,
    required this.postID,
    required this.initialCount,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.read(loveCountControllerProvider(Tuple2(postID, true)));
    final loveCount = ref.watch(loveCountProvider(Tuple2(postID, true)));
    final hasChangedLoveCount =
        ref.read(hasChangedLoveCountProvider(Tuple2(postID, true)));

    final count = (hasChangedLoveCount) ? loveCount : initialCount;

    return GestureDetector(
      child: Row(
        children: [
          Text(
            StringUtils.stringifyCount(count),
            style: TextStyle(
              fontSize: 17,
              fontWeight: FontWeight.w900,
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 0, 0, 0),
            child: Text(
              count < 2 ? "Love" : "Loves",
              style: TextStyle(
                fontSize: 17,
                color: Theme.of(context).colorScheme.secondary,
                fontWeight: FontWeight.w300,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _CommentCount extends ConsumerWidget {
  final String postID;
  final int initialCount;

  const _CommentCount({
    Key? key,
    required this.postID,
    required this.initialCount,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final commentCount = ref.watch(commentCountProvider(Tuple2(postID, true)));
    final hasChangedCommentCount =
        ref.read(hasChangedCommentCountProvider(Tuple2(postID, true)));

    final count = (hasChangedCommentCount) ? commentCount : initialCount;

    return GestureDetector(
      child: Row(
        children: [
          Text(
            StringUtils.stringifyCount(count),
            style: TextStyle(
              fontSize: 17,
              fontWeight: FontWeight.w900,
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 0, 0, 0),
            child: Text(
              count < 2 ? "Comment" : "Comments",
              style: TextStyle(
                fontSize: 17,
                color: Theme.of(context).colorScheme.secondary,
                fontWeight: FontWeight.w300,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _ActionButtons extends StatelessWidget {
  final PostModel postData;
  const _ActionButtons({
    Key? key,
    required this.postData,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 4, horizontal: 30),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: LoveReactIconButton(
              iconSize: 24,
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: CommentReplyButton(
              iconSize: 24,
            ),
          ),
          IconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.retweet),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 24,
          ),
          IconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.paperPlane),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 24,
          ),
        ],
      ),
    );
  }
}
