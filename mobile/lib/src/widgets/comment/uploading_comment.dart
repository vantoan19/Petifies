// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/models/comment.dart';
import 'package:mobile/src/models/uploading_comment.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/comment/comment_body.dart';
import 'package:mobile/src/widgets/comment/comment_head.dart';
import 'package:mobile/src/widgets/images/image_card.dart';
import 'package:mobile/src/widgets/videos/video_player_card.dart';

class UploadingComment extends StatelessWidget {
  const UploadingComment({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AbsorbPointer(
      absorbing: true,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 20),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.max,
          children: [
            CommentUserAvatar(),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const CommentHead(
                  isUploadingComment: true,
                ),
                Padding(
                  padding: const EdgeInsets.fromLTRB(0, 10, 0, 10),
                  child: CommentBody(
                    isUploadingComment: true,
                    width: MediaQuery.of(context).size.width - 126,
                    onlyTextFontSize: 17,
                    normalFontSize: 16,
                  ),
                ),
              ],
            )
          ],
        ),
      ),
    );
  }
}
