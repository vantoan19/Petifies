// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/comment/controller/list_comments_controller.dart';
import 'package:mobile/src/models/comment.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/context_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/widgets/appbars/comment_details_appbar.dart';
import 'package:mobile/src/widgets/buttons/comment_button.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/comment/comment_body.dart';
import 'package:mobile/src/widgets/comment/comment_footer.dart';
import 'package:mobile/src/widgets/comment/comment_head.dart';
import 'package:mobile/src/widgets/comment/uploading_comment.dart';
import 'package:mobile/src/widgets/images/image_card.dart';
import 'package:mobile/src/widgets/info_viewer/comment_count_view.dart';
import 'package:mobile/src/widgets/info_viewer/love_count_view.dart';
import 'package:mobile/src/widgets/posts/post_body.dart';
import 'package:mobile/src/widgets/posts/post_footer.dart';
import 'package:mobile/src/widgets/posts/post_head.dart';
import 'package:mobile/src/widgets/stepper/step.dart';
import 'package:mobile/src/widgets/textfield/reply_text_field.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';
import 'package:mobile/src/widgets/videos/video_player_card.dart';

class CommentDetailsScreenArguments {
  final PostModel postData;
  final CommentModel commentData;
  final List<CommentModel> ancestorComments;
  final bool autoFocus;
  CommentDetailsScreenArguments({
    required this.postData,
    required this.commentData,
    required this.ancestorComments,
    required this.autoFocus,
  });
}

class CommentDetailScreen extends ConsumerWidget {
  final bool autoFocus;
  const CommentDetailScreen({
    required this.autoFocus,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final commentData = ref.watch(commentInfoProvider);
    final _ =
        ref.watch(listCommentsControllerProvider(parentID: commentData.id));
    final comments = ref.watch(commentsProvider(parentID: commentData.id));
    final newlyCreatedComments =
        ref.watch(newlyCreatedCommentsProvider(parentID: commentData.id));

    return Scaffold(
      appBar: CommentDetailsAppbar(),
      body: Column(
        children: [
          Expanded(
            child: ListView.separated(
              padding: EdgeInsets.only(bottom: 200),
              itemBuilder: (context, index) {
                if (index == 0) {
                  return _CommentDetailsAndAncestors();
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

class _CommentDetailsAndAncestors extends ConsumerWidget {
  const _CommentDetailsAndAncestors({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final ancestorComments = ref.read(ancestorCommentsProvider);

    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: Constants.horizontalScreenPadding,
      ),
      child: Column(
        children: [
          ProviderScope(
            overrides: [
              isPostContextProvider.overrideWithValue(true),
              commentInfoProvider
                  .overrideWith((ref) => throw UnimplementedError()),
              ancestorCommentsProvider
                  .overrideWith((ref) => throw UnimplementedError()),
            ],
            child: _RootPost(),
          ),
          ...ancestorComments.sublist(0, ancestorComments.length - 1).map(
                (comment) => ProviderScope(
                  overrides: [
                    commentInfoProvider.overrideWithValue(comment),
                    isPostContextProvider.overrideWithValue(false),
                  ],
                  child: const _AncestorComment(),
                ),
              ),
          _CommentDetails(),
        ],
      ),
    );
  }
}

class _RootPost extends ConsumerWidget {
  const _RootPost({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final width = MediaQuery.of(context).size.width - 82;
    final postData = ref.watch(postInfoProvider);

    final postWidget = Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          const PostHeadInfo(isUploadingPost: false),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: PostBody(
              isUploadingPost: false,
              width: width,
              onlyTextFontSize: 18,
              normalFontSize: 16,
              spaceBetweenTextAndMedia: 8,
            ),
          ),
          PostFooter(
            iconSize: 16,
            textSize: 12,
            spaceBetween: 6,
          )
        ],
      ),
    );

    return CustomStep(
      leftWidget: UserAvatar(
        userAvatar: postData.owner.userAvatar,
      ),
      rightColumnWidget: postWidget,
      spaceBetween: 14,
      showStepLine: true,
    );
  }
}

class _AncestorComment extends ConsumerWidget {
  const _AncestorComment({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final width = MediaQuery.of(context).size.width - 82;
    final commentData = ref.watch(commentInfoProvider);

    final commentWidget = Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const CommentHead(
            isUploadingComment: false,
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: CommentBody(
              isUploadingComment: false,
              width: width,
              onlyTextFontSize: 18,
              normalFontSize: 16,
            ),
          ),
          const CommentFooter(
            iconSize: 16,
            textSize: 12,
            spaceBetween: 6,
            maxWidthBetweenLoveAndComment: 120,
          ),
        ],
      ),
    );

    return CustomStep(
      leftWidget: UserAvatar(
        userAvatar: commentData.owner.userAvatar,
      ),
      rightColumnWidget: commentWidget,
      spaceBetween: 14,
      showStepLine: true,
    );
  }
}

class _CommentDetails extends ConsumerWidget {
  const _CommentDetails({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final owner = ref.watch(commentInfoProvider.select((value) => value.owner));
    final width = MediaQuery.of(context).size.width -
        Constants.horizontalScreenPadding * 2;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            CommentUserAvatar(
              userAvatar: owner.userAvatar,
            ),
            Expanded(child: CommentHead(isUploadingComment: false)),
          ],
        ),
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 24),
          child: CommentBody(
            isUploadingComment: false,
            width: width,
            onlyTextFontSize: 24,
            normalFontSize: 18,
          ),
        ),
        _CommentDetailsFooter(),
      ],
    );
  }
}

class _CommentDetailsFooter extends ConsumerWidget {
  const _CommentDetailsFooter({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 0, 0, 0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Likes & comments data
          Column(
            children: [
              Container(
                padding: const EdgeInsets.fromLTRB(0, 4, 0, 16),
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
                  constraints: BoxConstraints(maxWidth: 240, minWidth: 200),
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
              // Action buttons
              _ActionButtons(),
            ],
          ),
        ],
      ),
    );
    ;
  }
}

class _ActionButtons extends StatelessWidget {
  const _ActionButtons({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 18, horizontal: 10),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: LoveReactIconButton(
              iconSize: 20,
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8),
            child: CommentReplyButton(
              iconSize: 20,
            ),
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
