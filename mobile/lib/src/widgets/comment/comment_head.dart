// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/utils/stringutils.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';
import 'package:mobile/src/widgets/comment/uploading_comment.dart';

class CommentHead extends ConsumerWidget {
  final bool isUploadingComment;

  const CommentHead({
    Key? key,
    required this.isUploadingComment,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    BasicUserInfoModel owner;
    DateTime createdAt;

    if (isUploadingComment) {
      owner =
          ref.watch(uploadingCommentInfoProvider.select((info) => info.owner));
      createdAt = ref
          .watch(uploadingCommentInfoProvider.select((info) => info.createdAt));
    } else {
      owner = ref.watch(commentInfoProvider.select((info) => info.owner));
      createdAt =
          ref.watch(commentInfoProvider.select((info) => info.createdAt));
    }

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.max,
      children: [
        Row(
          children: [
            Text(
              owner.firstName + " " + owner.lastName,
              style: TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.w800,
              ),
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(8.0, 0, 0, 0),
              child: Text(
                StringUtils.stringifyTime(createdAt),
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w300,
                  color: Colors.grey,
                ),
              ),
            ),
          ],
        ),
        // More button
        NoPaddingIconButton(
          onPressed: () {},
          icon: Icon(Icons.more_horiz),
          padding: EdgeInsets.zero,
          constraints: BoxConstraints(
            minHeight: 10,
            minWidth: 10,
            maxHeight: 20,
          ),
          color: Theme.of(context).colorScheme.secondary,
        )
      ],
    );
  }
}
