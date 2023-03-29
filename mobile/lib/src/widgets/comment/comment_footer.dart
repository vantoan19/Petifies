// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/widgets/buttons/comment_button.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';

class CommentFooter extends ConsumerWidget {
  const CommentFooter({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext contextm, WidgetRef ref) {
    return Container(
      constraints: BoxConstraints(maxWidth: 120),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          // Love reaction
          LoveReactButton(
            iconSize: 16,
            textSize: 12,
            spaceBetween: 6,
          ),
          // Comments
          CommentButton(
            textSize: 12,
            iconSize: 16,
            spaceBetween: 6,
          ),
        ],
      ),
    );
  }
}
