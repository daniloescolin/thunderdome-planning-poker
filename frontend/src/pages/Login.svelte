<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

    export let router
    export let xfetch
    export let notifications
    export let eventTag
    export let battleId

    const authMethod = appConfig.AuthMethod

    let warriorEmail = ''
    let warriorPassword = ''

    let warriorResetEmail = ''
    let forgotPassword = false

    $: targetPage = battleId
        ? `${appRoutes.battle}/${battleId}`
        : appRoutes.battles

    function authWarrior(e) {
        e.preventDefault()
        const body = {
            warriorEmail,
            warriorPassword,
        }

        xfetch('/api/auth', { body })
            .then(res => res.json())
            .then(function(newWarrior) {
                warrior.create({
                    id: newWarrior.id,
                    name: newWarrior.name,
                    email: newWarrior.email,
                    rank: newWarrior.rank,
                    notificationsEnabled: newWarrior.notificationsEnabled,
                })

                eventTag('login', 'engagement', 'success', () => {
                    router.route(targetPage, true)
                })
            })
            .catch(function(error) {
                notifications.danger($_('pages.login.authError'))
                eventTag('login', 'engagement', 'failure')
            })
    }

    function toggleForgotPassword() {
        forgotPassword = !forgotPassword
        eventTag(
            'forgot_password_toggle',
            'engagement',
            `forgot: ${forgotPassword}`,
        )
    }

    function sendPasswordReset(e) {
        e.preventDefault()
        const body = {
            warriorEmail: warriorResetEmail,
        }

        xfetch('/api/auth/forgot-password', { body })
            .then(function() {
                notifications.success(
                    $_('pages.login.sendResetSuccess', {
                        values: { email: warriorResetEmail },
                    }),
                    2000,
                )
                forgotPassword = !forgotPassword
                eventTag('forgot_password', 'engagement', 'success')
            })
            .catch(function(error) {
                notifications.danger($_('pages.login.sendResetError'))
                eventTag('forgot_password', 'engagement', 'failure')
            })
    }

    $: loginDisabled = warriorEmail === '' || warriorPassword === ''
    $: resetDisabled = warriorResetEmail === ''
</script>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
            {#if !forgotPassword}
                <form
                    on:submit="{authWarrior}"
                    class="bg-white shadow-lg rounded p-6 mb-4"
                    name="authWarrior">
                    <div
                        class="font-bold text-xl md:text-2xl mb-2 md:mb-6
                        md:leading-tight text-center">
                        {$_('pages.login.title')}
                    </div>
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourEmail">
                            {$_('pages.login.fields.email.label')}
                        </label>
                        <input
                            bind:value="{warriorEmail}"
                            placeholder="{$_('pages.login.fields.email.placeholder')}"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourEmail"
                            name="yourEmail"
                            type="email"
                            required />
                    </div>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourPassword">
                            {$_('pages.login.fields.password.label')}
                        </label>
                        <input
                            bind:value="{warriorPassword}"
                            placeholder="{$_('pages.login.fields.password.placeholder')}"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourPassword"
                            name="yourPassword"
                            type="password"
                            required />
                    </div>

                    <div class="text-right">
                        {#if authMethod === 'normal'}
                            <button
                                type="button"
                                class="inline-block align-baseline font-bold
                                text-sm text-blue-500 hover:text-blue-800 mr-4"
                                on:click="{toggleForgotPassword}">
                                {$_('pages.login.fields.password.forgotLabel')}
                            </button>
                        {/if}
                        <SolidButton type="submit" disabled="{loginDisabled}">
                            {$_('pages.login.button')}
                        </SolidButton>
                    </div>
                </form>
            {/if}

            {#if forgotPassword}
                <form
                    on:submit="{sendPasswordReset}"
                    class="bg-white shadow-lg rounded p-6 mb-4"
                    name="resetPassword">
                    <div
                        class="font-bold text-xl md:text-2xl mb-2 md:mb-6
                        md:leading-tight text-center">
                        Forgot Password
                    </div>
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourResetEmail">
                            Email
                        </label>
                        <input
                            bind:value="{warriorResetEmail}"
                            placeholder="Enter your email"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourResetEmail"
                            name="yourResetEmail"
                            type="email"
                            required />
                    </div>

                    <div class="text-right">
                        <button
                            type="button"
                            class="inline-block align-baseline font-bold text-sm
                            text-blue-500 hover:text-blue-800 mr-4"
                            on:click="{toggleForgotPassword}">
                            Cancel
                        </button>
                        <SolidButton type="submit" disabled="{resetDisabled}">
                            Send Reset Email
                        </SolidButton>
                    </div>
                </form>
            {/if}
        </div>
    </div>
</PageLayout>
