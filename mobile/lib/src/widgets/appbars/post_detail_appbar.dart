// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/widgets/buttons/go_back_button.dart';

class PostDetailsAppBar extends StatelessWidget implements PreferredSizeWidget {
  const PostDetailsAppBar({
    Key? key,
  }) : super(key: key);

  @override
  Size get preferredSize => const Size.fromHeight(70);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(
        "Post",
        style: TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w800,
        ),
      ),
      centerTitle: true,
      leading: const GoBackButton(),
    );
  }
}
