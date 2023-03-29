// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/models/uploading_post.dart';
import 'package:mobile/src/widgets/media_view/media_view.dart';
import 'package:mobile/src/widgets/posts/post_body.dart';
import 'package:mobile/src/widgets/posts/post_head.dart';

final uploadingPostInfoProvider =
    Provider<UploadingPostModel>((ref) => throw UnimplementedError());

class UploadingPost extends StatelessWidget {
  const UploadingPost({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AbsorbPointer(
      absorbing: true,
      child: Column(
        children: [
          LinearProgressIndicator(),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: Column(
              children: [
                const PostHead(
                  isUploadingPost: true,
                ),
                const PostBody(
                  isUploadingPost: true,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
