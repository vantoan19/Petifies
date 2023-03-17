// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/models/post.dart';
import 'package:mobile/src/widgets/posts/post_body.dart';
import 'package:mobile/src/widgets/posts/post_footer.dart';
import 'package:mobile/src/widgets/posts/post_head.dart';

class Post extends StatelessWidget {
  final PostModel postData;

  const Post({
    Key? key,
    required this.postData,
  }) : super(key: key);

  String _stringifyCount(int count) {
    if (count < 10000) {
      return count.toString();
    } else if (count < 1000000) {
      String count_ = (count.toDouble() / 1000.0).toStringAsFixed(1);
      return '${count_}K';
    } else if (count < 1000000000) {
      String count_ = (count.toDouble() / 1000000.0).toStringAsFixed(1);
      return '${count_}M';
    } else {
      String count_ = (count.toDouble() / 1000000000.0).toStringAsFixed(1);
      return '${count_}B';
    }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 16),
      child: Column(
        children: [
          PostHead(
            userAvatar: postData.owner.userAvatar,
            userName: postData.owner.firstName + " " + postData.owner.lastName,
            activity: postData.postActivity,
            postTime: postData.postTime.toString(),
          ),
          PostBody(
            textContent: postData.textContent,
            images: postData.images,
            videos: postData.videos,
          ),
          PostFooter(
            loveCount: _stringifyCount(postData.loveCount),
            commentCount: _stringifyCount(postData.commentCount),
          ),
        ],
      ),
    );
  }
}
