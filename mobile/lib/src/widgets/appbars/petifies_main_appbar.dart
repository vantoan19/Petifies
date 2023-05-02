import 'package:flutter/material.dart';
import 'package:mobile/src/widgets/buttons/go_back_button.dart';

class PetifiesMainAppbar extends StatelessWidget
    implements PreferredSizeWidget {
  const PetifiesMainAppbar({super.key});

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(
        "Petifies",
        style: TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w800,
        ),
      ),
      centerTitle: true,
      leading: const GoBackButton(),
      actions: [
        IconButton(
          iconSize: 26,
          onPressed: () {},
          icon: const Icon(Icons.settings),
        ),
        IconButton(
          iconSize: 26,
          padding: EdgeInsets.symmetric(horizontal: 30),
          onPressed: () {},
          icon: const Icon(Icons.chat_outlined),
        )
      ],
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(70);
}
