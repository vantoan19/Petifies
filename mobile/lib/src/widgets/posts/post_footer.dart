// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

class PostFooter extends StatelessWidget {
  final String loveCount;
  final String commentCount;

  const PostFooter({
    Key? key,
    required this.loveCount,
    required this.commentCount,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Row(
          children: [
            Row(
              children: [
                IconButton(
                    onPressed: () {}, icon: Icon(FontAwesomeIcons.heart)),
                Text(
                  loveCount,
                  style: TextStyle(
                    fontSize: 12,
                  ),
                ),
              ],
            ),
            Row(
              children: [
                IconButton(onPressed: () {}, icon: Icon(Icons.comment)),
                Text(
                  commentCount,
                  style: TextStyle(
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ],
        ),
      ],
    );
  }
}
