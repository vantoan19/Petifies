// ignore_for_file: public_member_api_docs, sort_constructors_first

import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/providers/comment_providers.dart';
import 'package:mobile/src/widgets/images/image_card.dart';
import 'package:mobile/src/widgets/videos/video_player_card.dart';
import 'package:video_player/video_player.dart';

class CommentBody extends ConsumerWidget {
  final bool isUploadingComment;
  final double width;
  final double onlyTextFontSize;
  final double normalFontSize;

  const CommentBody({
    Key? key,
    required this.isUploadingComment,
    required this.width,
    required this.onlyTextFontSize,
    required this.normalFontSize,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    String? textContent;
    String? imageURL;
    String? videoURL;
    File? imageFile;
    VideoPlayerController? videoController;

    if (isUploadingComment) {
      textContent = ref.watch(
          uploadingCommentInfoProvider.select((info) => info.textContent));
      imageURL = null;
      videoURL = null;
      imageFile =
          ref.watch(uploadingCommentInfoProvider.select((info) => info.image));
      videoController =
          ref.watch(uploadingCommentInfoProvider.select((info) => info.video));
    } else {
      textContent =
          ref.watch(commentInfoProvider.select((info) => info.textContent));
      imageURL =
          ref.watch(commentInfoProvider.select((info) => info.image))?.uri ??
              null;
      videoURL =
          ref.watch(commentInfoProvider.select((info) => info.video))?.uri ??
              null;
      imageFile = null;
      videoController = null;
    }

    final onlyText = (imageFile == null &&
        imageURL == null &&
        videoController == null &&
        videoURL == null);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (textContent != null)
          Padding(
            padding: EdgeInsets.fromLTRB(
              0,
              0,
              0,
              onlyText ? 0 : 6,
            ),
            child: Text(
              textContent,
              style: TextStyle(
                fontSize: onlyText ? onlyTextFontSize : normalFontSize,
                color: Theme.of(context).colorScheme.secondary,
              ),
            ),
          ),
        (imageURL != null || imageFile != null)
            ? ImageCard(
                isRoundedTopLeft: true,
                isRoundedTopRight: true,
                isRoundedBottomLeft: true,
                isRoundedBottomRight: true,
                width: width,
                maxHeight: 300,
                imageUrl: imageURL ?? "",
                imageFile: imageFile,
                isClickable: !isUploadingComment,
              )
            : const SizedBox.shrink(),
        (videoURL != null || videoController != null)
            ? VideoPlayerCard(
                isRoundedTopLeft: true,
                isRoundedTopRight: true,
                isRoundedBottomLeft: true,
                isRoundedBottomRight: true,
                videoUrl: videoURL ?? "",
                controller: videoController,
                autoPlay: true,
                width: width,
                maxHeight: 300,
                playNextVideoCallback: () => {},
                isClickable: !isUploadingComment,
              )
            : SizedBox.shrink(),
      ],
    );
  }
}
