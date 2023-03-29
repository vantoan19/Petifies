// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/buttons/go_back_button.dart';

class CreateCommentAppBar extends StatelessWidget
    implements PreferredSizeWidget {
  final String targetOwner;
  final String targetType;
  final VoidCallback addPostAction;

  const CreateCommentAppBar({
    Key? key,
    required this.addPostAction,
    required this.targetOwner,
    required this.targetType,
  }) : super(key: key);

  @override
  Size get preferredSize => const Size.fromHeight(70);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(
        "Reply to ${targetOwner}'s ${targetType}",
        style: TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w800,
        ),
      ),
      centerTitle: true,
      leading: const GoBackButton(),
      actions: [
        IconButton(
          iconSize: 30,
          onPressed: addPostAction,
          icon: Icon(
            Icons.add,
            color: Themes.blueColor,
          ),
        )
      ],
    );
  }
}
