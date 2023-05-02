// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/comment/screens/comment_details_screen.dart';
import 'package:mobile/src/features/home/navigators/home_navigator.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/widgets/comment/comment_body.dart';
import 'package:mobile/src/widgets/comment/comment_footer.dart';
import 'package:mobile/src/widgets/comment/comment_head.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';

class Comment extends ConsumerWidget {
  const Comment({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final owner = ref.watch(commentInfoProvider.select((value) => value.owner));

    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 16, horizontal: Constants.horizontalScreenPadding),
      child: GestureDetector(
        behavior: HitTestBehavior.translucent,
        onTap: () => Navigator.pushNamed(
          context,
          HomeNavigatorRoutes.commentDetails,
          arguments: CommentDetailsScreenArguments(
            postData: ref.read(postInfoProvider),
            commentData: ref.read(commentInfoProvider),
            ancestorComments: ref.read(commentInfoProvider).isPostParent
                ? [ref.read(commentInfoProvider)]
                : [
                    ...ref.read(ancestorCommentsProvider),
                    ref.read(commentInfoProvider),
                  ],
            autoFocus: false,
          ),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.max,
          children: [
            CommentUserAvatar(
              userAvatar: owner.userAvatar,
            ),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const CommentHead(
                    isUploadingComment: false,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 6, 0, 6),
                    child: CommentBody(
                      isUploadingComment: false,
                      width: MediaQuery.of(context).size.width - 122,
                      onlyTextFontSize: 17,
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
            )
          ],
        ),
      ),
    );
  }
}

class CommentUserAvatar extends UserAvatar {
  final String? userAvatar;
  const CommentUserAvatar({
    Key? key,
    this.userAvatar = null,
  }) : super(
          key: key,
          userAvatar: userAvatar,
          padding: const EdgeInsets.fromLTRB(0, 0, 14, 0),
        );
}
