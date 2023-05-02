// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/buttons/comment_button.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';

class PostFooter extends ConsumerWidget {
  final double iconSize;
  final double textSize;
  final double spaceBetween;
  const PostFooter({
    Key? key,
    required this.iconSize,
    required this.textSize,
    required this.spaceBetween,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Row(
      mainAxisSize: MainAxisSize.max,
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Flexible(
          child: LoveReactButton(
            textSize: textSize,
            iconSize: iconSize,
            spaceBetween: spaceBetween,
          ),
        ),
        // Comments
        Flexible(
          child: CommentButton(
            textSize: textSize,
            iconSize: iconSize,
            spaceBetween: spaceBetween,
          ),
        ),
        Flexible(
          child: NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.retweet),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: iconSize,
          ),
        ),
        Flexible(
          child: NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(Icons.bookmark_outline),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: iconSize + 4,
          ),
        ),
        Flexible(
          child: NoPaddingIconButton(
            onPressed: () {},
            icon: Icon(FontAwesomeIcons.paperPlane),
            color: Theme.of(context).colorScheme.secondary,
            iconSize: iconSize,
          ),
        )
      ],
    );
  }
}
