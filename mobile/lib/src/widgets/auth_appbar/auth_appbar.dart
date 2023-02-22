import 'package:flutter/material.dart';
import 'package:introduction_screen/introduction_screen.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/utils/navigation.dart';

class AuthAppBar extends StatelessWidget implements PreferredSizeWidget {
  final GlobalKey<IntroductionScreenState>? _introScreenKey;
  const AuthAppBar(
      {super.key, GlobalKey<IntroductionScreenState>? introScreenKey})
      : _introScreenKey = introScreenKey;

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Image.asset(
        Theme.of(context).brightness == Brightness.light
            ? Constants.logoLightThemePath
            : Constants.logoDarkThemePath,
        height: 60,
      ),
      centerTitle: true,
      leading: IconButton(
          onPressed: () => {
                if (_introScreenKey != null)
                  {
                    _introScreenKey?.currentState?.previous()
                  } // In case this is the introduction screen
                else
                  {NavigatorUtil.goBack(context)}
              },
          icon: const Icon(Icons.arrow_back)),
    );
  }
}
