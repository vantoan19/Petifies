import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/constants/languages.dart';
import 'package:mobile/src/features/auth/controllers/user_controllers.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/appbars/auth_appbar.dart';
import 'package:mobile/src/widgets/buttons/auth_button.dart';
import 'package:mobile/src/widgets/textfield/textformfield.dart';

class SignUpFormScreen extends StatelessWidget {
  const SignUpFormScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const AuthAppBar(),
      body: SingleChildScrollView(child: SignUpFormBody()),
    );
  }
}

class SignUpFormBody extends ConsumerStatefulWidget {
  const SignUpFormBody({
    super.key,
  });

  @override
  ConsumerState<ConsumerStatefulWidget> createState() => _SignUpFormBodyState();
}

class _SignUpFormBodyState extends ConsumerState<SignUpFormBody> {
  final _formKey = GlobalKey<FormState>();
  late String _email;
  late String _password;
  late String _firstName;
  late String _lastName;

  @override
  Widget build(BuildContext context) {
    ref.listen<AsyncValue<UserModel?>>(userCreateControllerProvider,
        (previous, next) {
      next.maybeWhen(
        data: (data) {
          if (data == null) return;
          NavigatorUtil.toSignIn(context);
        },
        orElse: () {},
      );
    });

    final createResult = ref.watch(userCreateControllerProvider);

    final errorMsg = createResult.maybeWhen(
      error: (error, stackTrace) => error.toString(),
      orElse: () => "nothing",
    );

    final isLoading = createResult.maybeWhen(
      data: (_) => createResult.isRefreshing,
      loading: () => true,
      orElse: () => false,
    );

    return Column(
      children: [
        const _GreetingHeader(),
        Text(errorMsg),
        Form(
          key: _formKey,
          autovalidateMode: AutovalidateMode.onUserInteraction,
          child: Column(
            children: [
              CustomTextFormField(
                label: Language.translate(context).emailLabel,
                icon: const Icon(Icons.email),
                onChange: (value) => _email = value,
              ),
              CustomTextFormField(
                label: Language.translate(context).passwordLabel,
                icon: const Icon(Icons.password),
                isObscureText: true,
                onChange: (value) => _password = value,
              ),
              CustomTextFormField(
                label: Language.translate(context).firstNameLabel,
                icon: const Icon(FontAwesomeIcons.signature),
                onChange: (value) => _firstName = value,
              ),
              CustomTextFormField(
                label: Language.translate(context).familyNameLabel,
                icon: const Icon(FontAwesomeIcons.peopleRoof),
                onChange: (value) => _lastName = value,
              ),
              SizedBox(
                height: 50,
              ),
              AuthButton(
                label: Language.translate(context).createAccountLabel,
                action: isLoading
                    ? null
                    : () {
                        if (!_formKey.currentState!.validate()) return;
                        ref.read(userCreateControllerProvider.notifier).handle(
                            email: _email,
                            password: _password,
                            firstName: _firstName,
                            lastName: _lastName);
                      },
                isLoading: isLoading,
              )
            ],
          ),
        ),
      ],
    );
  }
}

class _GreetingHeader extends StatelessWidget {
  const _GreetingHeader({super.key});

  @override
  Widget build(BuildContext context) {
    final TextStyle headerStyle = TextStyle(
        color: Theme.of(context).colorScheme.inversePrimary,
        fontSize: 32,
        fontWeight: FontWeight.w600);

    return Padding(
      padding: const EdgeInsets.symmetric(
          horizontal: Constants.horizontalScreenPadding, vertical: 50),
      child: Text(
        Language.translate(context).greetingRegistration,
        style: headerStyle,
      ),
    );
  }
}
