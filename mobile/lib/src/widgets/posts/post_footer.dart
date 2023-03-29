// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/widgets/buttons/comment_button.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';

class PostFooter extends ConsumerWidget {
  const PostFooter({
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(
        Constants.horizontalScreenPadding + 4,
        0,
        Constants.horizontalScreenPadding,
        0,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            constraints: BoxConstraints(maxWidth: 140, minWidth: 120),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              mainAxisSize: MainAxisSize.min,
              children: [
                // Love reaction
                LoveReactButton(
                  textSize: 14,
                  iconSize: 24,
                  spaceBetween: 12,
                ),
                // Comments
                CommentButton(
                  textSize: 14,
                  iconSize: 24,
                  spaceBetween: 12,
                )
              ],
            ),
          ),
          IconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.paperPlane),
            color: Theme.of(context).colorScheme.secondary,
          )
        ],
      ),
    );
  }
}
