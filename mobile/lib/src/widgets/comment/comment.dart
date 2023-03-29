// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/models/comment.dart';
import 'package:mobile/src/models/uploading_comment.dart';
import 'package:mobile/src/widgets/comment/comment_body.dart';
import 'package:mobile/src/widgets/comment/comment_footer.dart';
import 'package:mobile/src/widgets/comment/comment_head.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

final commentInfoProvider =
    Provider<CommentModel>((ref) => throw UnimplementedError());

class Comment extends ConsumerWidget {
  const Comment({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final owner = ref.watch(commentInfoProvider.select((value) => value.owner));

    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 8, horizontal: Constants.horizontalScreenPadding),
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
                  padding: const EdgeInsets.fromLTRB(0, 8, 0, 8),
                  child: const CommentBody(
                    isUploadingComment: false,
                  ),
                ),
                const CommentFooter(),
              ],
            ),
          )
        ],
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
