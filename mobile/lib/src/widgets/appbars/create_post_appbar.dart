// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:introduction_screen/introduction_screen.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/theme/themes.dart';
import 'package:mobile/src/utils/navigation.dart';

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
      leading: IconButton(
        onPressed: () => {NavigatorUtil.goBack(context)},
        icon: const Icon(Icons.arrow_back),
      ),
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
