// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/widgets/buttons/go_back_button.dart';

class OnlyGoBackAppbar extends StatelessWidget implements PreferredSizeWidget {
  const OnlyGoBackAppbar({
    Key? key,
  }) : super(key: key);

  @override
  Size get preferredSize => const Size.fromHeight(70);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      leading: const GoBackButton(),
    );
  }
}
