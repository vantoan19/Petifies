// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/widgets/media_view/media_view.dart';
import 'package:mobile/src/widgets/posts/post.dart';
import 'package:mobile/src/widgets/posts/uploading_post.dart';
import 'package:video_player/video_player.dart';

class PostBody extends ConsumerWidget {
  final bool isUploadingPost;

  const PostBody({
    required this.isUploadingPost,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    String? textContent;
    List<String>? imageURLs;
    List<String>? videoURLs;
    List<File>? imageFiles;
    List<VideoPlayerController>? videoControllers;

    if (isUploadingPost) {
      textContent = ref
          .watch(uploadingPostInfoProvider.select((info) => info.textContent));
      imageURLs = null;
      videoURLs = null;
      imageFiles =
          ref.watch(uploadingPostInfoProvider.select((info) => info.images));
      videoControllers =
          ref.watch(uploadingPostInfoProvider.select((info) => info.videos));
    } else {
      textContent =
          ref.watch(postInfoProvider.select((info) => info.textContent));
      imageURLs = ref
              .watch(postInfoProvider.select((info) => info.images))
              ?.map((e) => e.uri)
              .toList() ??
          null;
      videoURLs = ref
              .watch(postInfoProvider.select((info) => info.videos))
              ?.map((e) => e.uri)
              .toList() ??
          null;
      imageFiles = null;
      videoControllers = null;
    }

    bool onlyText = (imageURLs == null || imageURLs.length == 0) &
        (videoURLs == null || videoURLs.length == 0) &
        (imageFiles == null || imageFiles.length == 0) &
        (videoControllers == null || videoControllers.length == 0);

    bool onlyMedia = (textContent == null);

    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 20, 0, 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Text Content
          if (textContent != null && textContent != "")
            Padding(
              padding: EdgeInsets.fromLTRB(
                Constants.horizontalScreenPadding + 4,
                0,
                Constants.horizontalScreenPadding,
                onlyText ? 0 : 16,
              ),
              child: Align(
                child: Text(
                  textContent,
                  style: TextStyle(
                    fontSize: onlyText ? 24 : 18,
                  ),
                ),
                alignment: Alignment.topLeft,
              ),
            ),
          // Image & Video content
          if ((imageURLs != null && imageURLs.length > 0) ||
              (videoURLs != null && videoURLs.length > 0) ||
              (imageFiles != null && imageFiles.length > 0) ||
              (videoControllers != null && videoControllers.length > 0))
            MediaView(
              imageUrls: imageURLs != null ? imageURLs : [],
              videoUrls: videoURLs != null ? videoURLs : [],
              imageFiles: imageFiles,
              videoControllers: videoControllers,
              isClickable: !isUploadingPost,
            ),
        ],
      ),
    );
  }
}
