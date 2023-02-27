import 'package:flutter/material.dart';
import 'package:mobile/src/constants/constants.dart';

class MainAppBar extends StatelessWidget implements PreferredSizeWidget {
  const MainAppBar({super.key});

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Row(
        children: [
          SizedBox(
            width: 10,
          ),
          Image.asset(
            Theme.of(context).brightness == Brightness.light
                ? Constants.logoLightThemePath
                : Constants.logoDarkThemePath,
            height: 60,
          ),
        ],
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.start,
      ),
      automaticallyImplyLeading: false,
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
