import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/custom_widgets/auth_appbar/auth_appbar.dart';
import 'package:mobile/src/custom_widgets/auth_button/auth_button.dart';
import 'package:mobile/src/custom_widgets/auth_textfield/auth_textfiled.dart';

class SignInScreen extends StatefulWidget {
  const SignInScreen({super.key});

  @override
  State<SignInScreen> createState() => _SignInScreenState();
}

class _SignInScreenState extends State<SignInScreen> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const AuthAppBar(),
      body: SingleChildScrollView(
        child: _SignInBody(
          emailController: _emailController,
          passwordController: _passwordController,
        ),
      ),
    );
  }
}

class _SignInBody extends StatelessWidget {
  final TextEditingController _emailController;
  final TextEditingController _passwordController;

  const _SignInBody(
      {super.key, required emailController, required passwordController})
      : _emailController = emailController,
        _passwordController = passwordController;

  @override
  Widget build(BuildContext context) {
    double width = MediaQuery.of(context).size.width;

    return Column(mainAxisAlignment: MainAxisAlignment.start, children: [
      SvgPicture.asset(
        Constants.signInImagePath,
        width: width,
      ),
      AuthTextField(
          label: 'E-Mail',
          icon: const Icon(Icons.email),
          controller: _emailController),
      AuthTextField(
        label: 'Password',
        icon: const Icon(Icons.lock),
        controller: _passwordController,
        isObscureText: true,
      ),
      const _SignInForgotPasswordButton(),
      AuthButton(label: 'Sign In', action: () => {}),
      ThirdpartyAuthButton.withColor(
        label: 'Sign in with Google',
        action: () => {},
        icon: const FaIcon(FontAwesomeIcons.google),
        color: Theme.of(context).colorScheme.secondary,
      ),
      const _SignInToRegisterButton()
    ]);
  }
}

class _SignInForgotPasswordButton extends StatelessWidget {
  const _SignInForgotPasswordButton({super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
        horizontal: Constants.horizontalScreenPadding,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          TextButton(
              onPressed: () => {}, child: const Text('Forgot your password')),
        ],
      ),
    );
  }
}

class _SignInToRegisterButton extends StatelessWidget {
  const _SignInToRegisterButton({super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        Row(
          children: [
            const Text('New to Petifies?'),
            TextButton(
                onPressed: () => {Navigator.of(context).pushNamed('/signup')},
                child: const Text('Register'))
          ],
        )
      ],
    );
  }
}
