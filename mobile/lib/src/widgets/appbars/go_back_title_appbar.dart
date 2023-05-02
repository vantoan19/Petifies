// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/widgets/buttons/go_back_button.dart';

class GoBackTitleAppbar extends StatelessWidget implements PreferredSizeWidget {
  final String title;
  const GoBackTitleAppbar({
    Key? key,
    required this.title,
  }) : super(key: key);

  @override
  Size get preferredSize => const Size.fromHeight(70);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(
        title,
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
