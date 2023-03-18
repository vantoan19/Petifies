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
    return Padding(
      padding: const EdgeInsets.fromLTRB(12, 0, 8, 0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
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
              SizedBox(
                width: 10,
              ),
              Row(
                children: [
                  IconButton(
                    onPressed: () {},
                    icon: Icon(Icons.comment_outlined),
                  ),
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
          IconButton(onPressed: () {}, icon: Icon(Icons.send_outlined))
        ],
      ),
    );
  }
}