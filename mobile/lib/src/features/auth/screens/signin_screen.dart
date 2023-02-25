import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/main.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/constants/languages.dart';
import 'package:mobile/src/features/auth/controllers/auth_controllers.dart';
import 'package:mobile/src/providers/model_providers.dart';
import 'package:mobile/src/providers/secure_storage_provider.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/auth_appbar/auth_appbar.dart';
import 'package:mobile/src/widgets/textfield/textformfield.dart';
import 'package:mobile/src/widgets/buttons/auth_button.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class SignInScreen extends StatefulWidget {
  const SignInScreen({super.key});

  @override
  State<SignInScreen> createState() => _SignInScreenState();
}

class _SignInScreenState extends State<SignInScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const AuthAppBar(),
      body: SingleChildScrollView(
        child: _SignInBody(),
      ),
    );
  }
}

class _SignInBody extends ConsumerStatefulWidget {
  const _SignInBody({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() => _SignInBodyState();
}

class _SignInBodyState extends ConsumerState<_SignInBody> {
  final _formKey = GlobalKey<FormState>();
  late String _email;
  late String _password;

  @override
  Widget build(BuildContext context) {
    final double width = MediaQuery.of(context).size.width;
    ref.listen<AsyncValue<Map<String, dynamic>?>>(loginControllerProvider,
        (previous, next) {
      next.maybeWhen(
        data: (data) async {
          if (data == null) return;
          final secureStorage = ref.read(secureStorageProvider);
          await secureStorage.write(
              key: "accessToken", value: data["tokens"].accessToken);
          await secureStorage.write(
              key: "accessTokenExpiresAt",
              value: data["tokens"].accessTokenExpiresAt.toString());
          await secureStorage.write(
              key: "refreshToken", value: data["tokens"].refreshToken);
          await secureStorage.write(
              key: "refreshTokenExpiresAt",
              value: data["tokens"].refreshTokenExpiresAt.toString());

          ref.read(myUserProvider.notifier).SetUser(data["user"]);
          NavigatorUtil.toHomePage(context);
        },
        orElse: () {},
      );
    });

    final loginResult = ref.watch(loginControllerProvider);

    final errorMsg = loginResult.maybeWhen(
      error: (error, stackTrace) => error.toString(),
      orElse: () => "nothing",
    );

    final isLoading = loginResult.maybeWhen(
      data: (_) => loginResult.isRefreshing,
      loading: () => true,
      orElse: () => false,
    );

    return Column(mainAxisAlignment: MainAxisAlignment.start, children: [
      SvgPicture.asset(
        Constants.signInImagePath,
        width: width,
      ),
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
              icon: const Icon(Icons.lock),
              isObscureText: true,
              onChange: (value) => _password = value,
            ),
            const _ForgotPasswordButton(),
            AuthButton(
              label: Language.translate(context).signInLabel,
              isLoading: isLoading,
              action: () {
                if (!_formKey.currentState!.validate()) return;
                ref
                    .read(loginControllerProvider.notifier)
                    .handle(email: _email, password: _password);
              },
            ),
          ],
        ),
      ),
      const _SignInWithGoogleButoon(),
      const _RegisterButton()
    ]);
  }
}

class _ForgotPasswordButton extends StatelessWidget {
  const _ForgotPasswordButton({super.key});

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
            onPressed: () => {},
            child: Text(Language.translate(context).forgotPasswordText),
          ),
        ],
      ),
    );
  }
}

class _SignInWithGoogleButoon extends StatelessWidget {
  const _SignInWithGoogleButoon({super.key});

  @override
  Widget build(BuildContext context) {
    return ThirdpartyAuthButton.withColor(
      label: Language.translate(context).googleSignInLabel,
      action: () => {},
      icon: const FaIcon(FontAwesomeIcons.google),
      color: Theme.of(context).colorScheme.secondary,
    );
  }
}

class _RegisterButton extends StatelessWidget {
  const _RegisterButton({super.key});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        Row(
          children: [
            Text(Language.translate(context).newbieText),
            TextButton(
                onPressed: () => {Navigator.of(context).pushNamed('/signup')},
                child: Text(Language.translate(context).registerLabel))
          ],
        )
      ],
    );
  }
}
