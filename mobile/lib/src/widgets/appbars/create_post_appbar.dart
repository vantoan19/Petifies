// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/widgets/buttons/go_back_button.dart';

class CreatePostAppBar extends StatelessWidget implements PreferredSizeWidget {
  final VoidCallback addPostAction;
  const CreatePostAppBar({
    Key? key,
    required this.addPostAction,
  }) : super(key: key);

  @override
  Size get preferredSize => const Size.fromHeight(70);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(
        "New post",
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
