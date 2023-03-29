// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/features/media/screens/media_full_page_screen.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/posts/post.dart';

class ImageCard extends ConsumerWidget {
  final bool isRoundedTopLeft;
  final bool isRoundedTopRight;
  final bool isRoundedBottomLeft;
  final bool isRoundedBottomRight;
  final double? height;
  final double? width;
  final double maxHeight;
  final EdgeInsetsGeometry? padding;
  final String imageUrl;
  final File? imageFile;

  final bool isClickable;

  const ImageCard({
    Key? key,
    required this.isRoundedTopLeft,
    required this.isRoundedTopRight,
    required this.isRoundedBottomLeft,
    required this.isRoundedBottomRight,
    this.height = null,
    this.width = null,
    required this.maxHeight,
    this.padding = null,
    required this.imageUrl,
    this.imageFile = null,
    required this.isClickable,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final image = (imageFile != null)
        ? Image.file(
            imageFile!,
            fit: BoxFit.cover,
          )
        : Image.network(
            imageUrl,
            loadingBuilder: (context, child, loadingProgress) =>
                (loadingProgress == null) ? child : CircularProgressIndicator(),
            fit: BoxFit.cover,
          );

    Widget imageCard = Container(
      width: width,
      height: height,
      padding: padding,
      constraints: BoxConstraints(
        maxHeight: maxHeight,
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.only(
          topLeft: isRoundedTopLeft ? Radius.circular(12) : Radius.zero,
          topRight: isRoundedTopRight ? Radius.circular(12) : Radius.zero,
          bottomLeft: isRoundedBottomLeft ? Radius.circular(12) : Radius.zero,
          bottomRight: isRoundedBottomRight ? Radius.circular(12) : Radius.zero,
        ),
        child: Container(
          child: image,
        ),
      ),
    );

    if (!isClickable) {
      return imageCard;
    }

    return GestureDetector(
        onTap: () {
          showModalBottomSheet(
              context: context,
              isScrollControlled: true,
              useSafeArea: true,
              barrierColor: Theme.of(context).scaffoldBackgroundColor,
              builder: (context) {
                final isPostTarget = ref.read(isPostContextProvider);
                return ProviderScope(
                  overrides: [
                    isPostContextProvider.overrideWithValue(isPostTarget),
                    postInfoProvider
                        .overrideWithValue(ref.read(postInfoProvider)),
                    if (!isPostTarget)
                      commentInfoProvider
                          .overrideWithValue(ref.read(commentInfoProvider)),
                  ],
                  child: MediaFullPageScreen(
                    mediaUrl: imageUrl,
                    isMediaImage: true,
                  ),
                );
              });
        },
        child: imageCard);
  }
}
