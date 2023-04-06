// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/home/navigators/home_navigator.dart';
import 'package:mobile/src/features/post/screens/post_detail_screen.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/providers/post_providers.dart';
import 'package:mobile/src/widgets/posts/post_body.dart';
import 'package:mobile/src/widgets/posts/post_footer.dart';
import 'package:mobile/src/widgets/posts/post_head.dart';

class Post extends ConsumerWidget {
  const Post({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Padding(
      padding: const EdgeInsets.symmetric(
        vertical: 12,
        horizontal: Constants.horizontalScreenPadding,
      ),
      child: GestureDetector(
        onTap: () => Navigator.pushNamed(
          context,
          HomeNavigatorRoutes.postDetails,
          arguments: PostDetailsScreenArguments(
            postData: ref.read(postInfoProvider),
            autoFocus: false,
          ),
        ),
        behavior: HitTestBehavior.translucent,
        child: Column(
          children: [
            const PostHead(
              isUploadingPost: false,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 12),
              child: PostBody(
                isUploadingPost: false,
                width: MediaQuery.of(context).size.width -
                    Constants.horizontalScreenPadding * 2,
                onlyTextFontSize: 24,
                normalFontSize: 18,
                spaceBetweenTextAndMedia: 8,
              ),
            ),
            const PostFooter(
              iconSize: 20,
              textSize: 14,
              spaceBetween: 6,
            ),
          ],
        ),
      ),
    );
  }
}
