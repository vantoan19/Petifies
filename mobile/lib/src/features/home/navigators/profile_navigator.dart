// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:mobile/src/features/profile/screens/my_profile_screen.dart';
import 'package:mobile/src/features/profile/screens/profile_screen.dart';

class ProfileNavigatorRoutes {
  static const String root = '/';
  static const String myProfile = '/my-profile';
}

class ProfileNavigator extends StatelessWidget {
  final GlobalKey<NavigatorState> navigatorKey;

  const ProfileNavigator({
    Key? key,
    required this.navigatorKey,
  }) : super(key: key);

  void _pushMyProfile(BuildContext context) {
    Navigator.pushNamed(
      context,
      ProfileNavigatorRoutes.myProfile,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Navigator(
      key: navigatorKey,
      initialRoute: ProfileNavigatorRoutes.root,
      onGenerateRoute: onGenerateRoute,
    );
  }

  Route onGenerateRoute(RouteSettings settings) {
    switch (settings.name) {
      case ProfileNavigatorRoutes.root:
        return MaterialPageRoute(
          builder: (context) =>
              ProfileScreen(toMyProfileCallback: () => _pushMyProfile(context)),
        );
      case ProfileNavigatorRoutes.myProfile:
        return MaterialPageRoute(builder: (context) => MyProfileScreen());
      default:
        return MaterialPageRoute(
          builder: (context) =>
              ProfileScreen(toMyProfileCallback: () => _pushMyProfile(context)),
        );
    }
  }
}
