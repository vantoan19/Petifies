// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/profile/screens/my_profile_screen.dart';
import 'package:mobile/src/providers/model_providers.dart';

class ProfileScreen extends StatelessWidget {
  final Function(PreferredSizeWidget?, Widget) navigateCallback;

  const ProfileScreen({
    Key? key,
    required this.navigateCallback,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: [
        _UserProfile(
          navigateCallback: navigateCallback,
        ),
        Divider(
          thickness: 5,
          color: Color.fromRGBO(10, 10, 10, 10),
        ),
        Padding(
          padding: const EdgeInsets.all(25.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _MenuButton(
                label: "Petifies",
                action: () {},
                icon: Icon(Icons.pets),
              ),
              _MenuButton(
                label: "New Feeds",
                action: () {},
                icon: Icon(Icons.feed),
              ),
              _MenuButton(
                label: "Stories",
                action: () {},
                icon: Icon(Icons.tv),
              ),
              _MenuButton(
                label: "Friends",
                action: () {},
                icon: Icon(Icons.group),
              ),
            ],
          ),
        ),
        Divider(
          thickness: 5,
          color: Color.fromRGBO(10, 10, 10, 10),
        ),
        Padding(
          padding: const EdgeInsets.all(13.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsets.all(12.0),
                child: Text(
                  "Account settings",
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w400,
                  ),
                ),
              ),
              _AccountSettingButton(
                label: "Personal information",
                action: () {},
                icon: Icon(Icons.person),
              ),
              _AccountSettingButton(
                label: "Password and security",
                action: () {},
                icon: Icon(Icons.security),
              ),
              _AccountSettingButton(
                label: "Languages",
                action: () {},
                icon: Icon(Icons.language),
              ),
              _AccountSettingButton(
                label: "Notifications",
                action: () {},
                icon: Icon(Icons.notifications),
              ),
              _AccountSettingButton(
                label: "Privacy and sharing",
                action: () {},
                icon: Icon(Icons.lock),
              ),
              _AccountSettingButton(
                label: "Posts and stories",
                action: () {},
                icon: Icon(Icons.dynamic_feed),
              ),
              _AccountSettingButton(
                label: "Themes",
                action: () {},
                icon: Icon(Icons.dark_mode),
              ),
            ],
          ),
        )
      ],
    );
  }
}

class _UserProfile extends ConsumerWidget {
  final Function(PreferredSizeWidget?, Widget) navigateCallback;

  const _UserProfile({
    required this.navigateCallback,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(myUserProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 25, horizontal: 35),
      child: GestureDetector(
        behavior: HitTestBehavior.translucent,
        onTap: () {
          navigateCallback(null, MyProfileScreen());
        },
        child: Container(
          child: Row(
            children: [
              (userInfo?.userAvatar != null)
                  ? CircleAvatar(
                      backgroundImage: NetworkImage(
                        userInfo!.userAvatar!,
                      ),
                      backgroundColor: Colors.transparent,
                      radius: 35,
                    )
                  : CircleAvatar(
                      backgroundImage: AssetImage(
                        Constants.defaultAvatarPng,
                      ),
                      backgroundColor: Colors.transparent,
                      radius: 35,
                    ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 0, 0, 0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      userInfo!.firstName + " " + userInfo.lastName,
                      style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const Text(
                      "see your profile",
                      style: TextStyle(
                        color: Colors.grey,
                        fontSize: 15,
                        fontWeight: FontWeight.w200,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _MenuButton extends StatelessWidget {
  final String label;
  final VoidCallback action;
  final Widget icon;

  const _MenuButton({
    Key? key,
    required this.label,
    required this.action,
    required this.icon,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(6),
      child: GestureDetector(
        onTap: action,
        child: Row(
          children: [
            icon,
            Padding(
              padding: const EdgeInsets.fromLTRB(25, 0, 0, 0),
              child: Text(
                label,
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.w700,
                ),
              ),
            )
          ],
        ),
      ),
    );
  }
}

class _AccountSettingButton extends StatelessWidget {
  final String label;
  final VoidCallback action;
  final Widget icon;

  const _AccountSettingButton({
    Key? key,
    required this.label,
    required this.action,
    required this.icon,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        border: Border(
          bottom: BorderSide(
              color: Theme.of(context).colorScheme.secondary, width: 0.3),
        ),
      ),
      child: GestureDetector(
        onTap: action,
        child: Row(
          children: [
            icon,
            Padding(
              padding: const EdgeInsets.fromLTRB(25, 0, 0, 0),
              child: Text(
                label,
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.w200,
                ),
              ),
            )
          ],
        ),
      ),
    );
  }
}
