// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/widgets/posts/post_body.dart';
import 'package:mobile/src/widgets/posts/post_footer.dart';
import 'package:mobile/src/widgets/posts/post_head.dart';

final postInfoProvider =
    Provider<PostModel>((ref) => throw UnimplementedError());

class Post extends ConsumerWidget {
  const Post({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 16),
      child: GestureDetector(
        onTap: () => Navigator.pushNamed(
          context,
          '/post-details',
          arguments: Tuple2(ref.read(postInfoProvider), false),
        ),
        behavior: HitTestBehavior.translucent,
        child: Column(
          children: [
            const PostHead(
              isUploadingPost: false,
            ),
            const PostBody(
              isUploadingPost: false,
            ),
            const PostFooter(),
          ],
        ),
      ),
    );
  }
}
