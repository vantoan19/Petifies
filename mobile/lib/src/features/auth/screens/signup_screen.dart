import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/custom_widgets/auth_appbar/auth_appbar.dart';
import 'package:mobile/src/custom_widgets/auth_button/auth_button.dart';
import 'package:mobile/src/custom_widgets/text_divider/text_divider.dart';
import 'package:mobile/src/theme/themes.dart';

class SignUpScreen extends StatelessWidget {
  const SignUpScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      appBar: AuthAppBar(),
      body: SingleChildScrollView(child: _SignUpBody()),
    );
  }
}

class _SignUpBody extends StatelessWidget {
  const _SignUpBody({super.key});

  @override
  Widget build(BuildContext context) {
    double width = MediaQuery.of(context).size.width;

    return Column(
      children: [
        Container(
          height: 20,
        ),
        SvgPicture.asset(
          Constants.signUpImagePath,
          width: width,
        ),
        const _SignUpButtons(),
        const _SignUpAgreement()
      ],
    );
  }
}

class _SignUpButtons extends StatelessWidget {
  const _SignUpButtons({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 80, 0, 40),
      child: Column(children: [
        AuthButton(
            label: 'Create account',
            action: () => {Navigator.of(context).pushNamed('/signup/form')}),
        const TextDivider(
          text: "or",
          thickness: 1.5,
        ),
        ThirdpartyAuthButton.withColor(
          icon: const Icon(FontAwesomeIcons.google),
          label: 'Create account with Google',
          action: () => {},
          color: Themes.greyColor,
        )
      ]),
    );
  }
}

class _SignUpAgreement extends StatelessWidget {
  const _SignUpAgreement({super.key});

  @override
  Widget build(BuildContext context) {
    TextStyle textStyle = const TextStyle(color: Themes.greyColor);
    TextStyle linkStyle = const TextStyle(color: Themes.blueColor);

    return Padding(
        padding: const EdgeInsets.symmetric(
            horizontal: Constants.horizontalScreenPadding),
        child: RichText(
          text: TextSpan(children: [
            TextSpan(
                text: "By signing up, you agree to our ", style: textStyle),
            TextSpan(
                text: "Terms of Service",
                style: linkStyle,
                recognizer: TapGestureRecognizer()..onTap = () {}),
            TextSpan(text: ", ", style: textStyle),
            TextSpan(
                text: "Privacy Policy",
                style: linkStyle,
                recognizer: TapGestureRecognizer()..onTap = () {}),
            TextSpan(text: " and ", style: textStyle),
            TextSpan(
                text: "Uses of Cookies",
                style: linkStyle,
                recognizer: TapGestureRecognizer()..onTap = () {}),
          ]),
        ));
  }
}
