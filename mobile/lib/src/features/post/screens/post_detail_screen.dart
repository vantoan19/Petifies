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
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/appbars/post_detail_appbar.dart';
import 'package:mobile/src/widgets/buttons/comment_button.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/comment/uploading_comment.dart';
import 'package:mobile/src/widgets/info_viewer/comment_count_view.dart';
import 'package:mobile/src/widgets/info_viewer/love_count_view.dart';
import 'package:mobile/src/widgets/media_view/media_view.dart';
import 'package:mobile/src/widgets/posts/post_body.dart';
import 'package:mobile/src/widgets/posts/post_head.dart';
import 'package:mobile/src/widgets/textfield/reply_text_field.dart';

class PostDetailsScreenArguments {
  final PostModel postData;
  final bool autoFocus;
  PostDetailsScreenArguments({
    required this.postData,
    required this.autoFocus,
  });
}

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
              padding: EdgeInsets.only(bottom: 200),
              itemBuilder: (context, index) {
                if (index == 0) {
                  return _PostDetails(
                    postData: postData,
                  );
                } else if (comments.length == 0) {
                  return SizedBox.shrink();
                } else {
                  if (newlyCreatedComments.length > 0 &&
                      index - 1 < newlyCreatedComments.length) {
                    return ProviderScope(
                      key: ObjectKey(newlyCreatedComments[index - 1].tempId),
                      overrides: [
                        uploadingCommentInfoProvider
                            .overrideWithValue(newlyCreatedComments[index - 1])
                      ],
                      child: const UploadingComment(),
                    );
                  }
                  return ProviderScope(
                    key: ObjectKey(
                        comments[index - newlyCreatedComments.length - 1].id),
                    overrides: [
                      commentInfoProvider.overrideWithValue(
                          comments[index - newlyCreatedComments.length - 1]),
                      isPostContextProvider.overrideWithValue(false),
                    ],
                    child: const Comment(),
                  );
                }
              },
              separatorBuilder: (context, index) {
                return Divider(
                  thickness: 0.15,
                  height: 1,
                  color: Theme.of(context).colorScheme.secondary,
                );
              },
              itemCount: (comments.length > 0 ? comments.length : 1) + 1,
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

class _PostDetails extends StatelessWidget {
  final PostModel postData;

  const _PostDetails({
    Key? key,
    required this.postData,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.of(context).size.width -
        2 * Constants.horizontalScreenPadding;

    return Padding(
      padding: const EdgeInsets.symmetric(
          horizontal: Constants.horizontalScreenPadding),
      child: Column(
        mainAxisSize: MainAxisSize.max,
        children: [
          const PostHead(
            isUploadingPost: false,
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(0, 24, 0, 0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                PostBody(
                  isUploadingPost: false,
                  width: width,
                  onlyTextFontSize: 24,
                  normalFontSize: 22,
                  spaceBetweenTextAndMedia: 12,
                ),
                // Count & actions
                Column(
                  children: [
                    Container(
                      padding: const EdgeInsets.fromLTRB(0, 28, 0, 16),
                      decoration: BoxDecoration(
                        border: Border(
                          bottom: BorderSide(
                            color: Theme.of(context).colorScheme.secondary,
                            width: 0.15,
                          ),
                        ),
                      ),
                      width: MediaQuery.of(context).size.width * 0.96,
                      alignment: Alignment.centerLeft,
                      child: Container(
                        constraints:
                            BoxConstraints(maxWidth: 240, minWidth: 200),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const LoveCountView(fontSize: 16),
                            const CommentCountView(fontSize: 16),
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
      padding: const EdgeInsets.symmetric(vertical: 18, horizontal: 10),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          LoveReactIconButton(
            iconSize: 20,
          ),
          CommentReplyButton(
            iconSize: 20,
          ),
          NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.retweet),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 20,
          ),
          NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(Icons.bookmark_outline),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 24,
          ),
          NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.paperPlane),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: 20,
          ),
        ],
      ),
    );
  }
}
