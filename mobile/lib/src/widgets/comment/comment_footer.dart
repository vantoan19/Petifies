// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/widgets/buttons/comment_button.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';

class CommentFooter extends ConsumerWidget {
  final double iconSize;
  final double textSize;
  final double spaceBetween;
  final double maxWidthBetweenLoveAndComment;

  const CommentFooter({
    Key? key,
    required this.iconSize,
    required this.textSize,
    required this.spaceBetween,
    required this.maxWidthBetweenLoveAndComment,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Row(
      mainAxisSize: MainAxisSize.max,
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Container(
          constraints: BoxConstraints(
            maxWidth: maxWidthBetweenLoveAndComment,
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              // Love reaction
              LoveReactButton(
                iconSize: iconSize,
                textSize: textSize,
                spaceBetween: spaceBetween,
              ),
              // Comments
              CommentButton(
                iconSize: iconSize,
                textSize: textSize,
                spaceBetween: spaceBetween,
              ),
            ],
          ),
        ),
        NoPaddingIconButton(
          onPressed: () {},
          icon: Icon(FontAwesomeIcons.paperPlane),
          color: Theme.of(context).colorScheme.secondary,
          iconSize: iconSize,
        )
      ],
    );
  }
}
