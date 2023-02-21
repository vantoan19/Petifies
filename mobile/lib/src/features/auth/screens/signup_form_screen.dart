import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/constants/languages.dart';
import 'package:mobile/src/custom_widgets/auth_appbar/auth_appbar.dart';
import 'package:mobile/src/custom_widgets/buttons/auth_button.dart';
import 'package:mobile/src/custom_widgets/auth_textfield/auth_textfiled.dart';

class SignUpFormScreen extends StatefulWidget {
  const SignUpFormScreen({super.key});

  @override
  State<SignUpFormScreen> createState() => _SignUpFormScreenState();
}

class _SignUpFormScreenState extends State<SignUpFormScreen> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _firstNameController = TextEditingController();
  final TextEditingController _familyNameController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const AuthAppBar(),
      body: SingleChildScrollView(
          child: SignUpFormBody(
        emaiController: _emailController,
        passwordController: _passwordController,
        firstNameController: _firstNameController,
        familyNameController: _familyNameController,
      )),
    );
  }
}

class SignUpFormBody extends StatelessWidget {
  final TextEditingController _emailController;
  final TextEditingController _passwordController;
  final TextEditingController _firstNameController;
  final TextEditingController _familyNameController;

  const SignUpFormBody(
      {super.key,
      required emaiController,
      required passwordController,
      required firstNameController,
      required familyNameController})
      : _emailController = emaiController,
        _passwordController = passwordController,
        _firstNameController = firstNameController,
        _familyNameController = familyNameController;

  @override
  Widget build(BuildContext context) {
    final TextStyle headerStyle = TextStyle(
        color: Theme.of(context).colorScheme.inversePrimary,
        fontSize: 32,
        fontWeight: FontWeight.w600);

    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(
              horizontal: Constants.horizontalScreenPadding, vertical: 50),
          child: Text(
            Language.translate(context).greetingRegistration,
            style: headerStyle,
          ),
        ),
        AuthTextField(
            label: Language.translate(context).emailLabel,
            icon: const Icon(Icons.email),
            controller: _emailController),
        AuthTextField(
          label: Language.translate(context).passwordLabel,
          icon: const Icon(Icons.password),
          controller: _passwordController,
          isObscureText: true,
        ),
        AuthTextField(
            label: Language.translate(context).firstNameLabel,
            icon: const Icon(FontAwesomeIcons.signature),
            controller: _firstNameController),
        AuthTextField(
            label: Language.translate(context).familyNameLabel,
            icon: const Icon(FontAwesomeIcons.peopleRoof),
            controller: _familyNameController),
        Container(
          height: 50,
        ),
        AuthButton(
          label: Language.translate(context).createAccountLabel,
          action: () => {},
        )
      ],
    );
  }
}
