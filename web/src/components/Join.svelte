<script lang="ts">
  import { fly } from 'svelte/transition';
  import { api, type Decision, type Me } from '../lib/api';
  import { lastName } from '../lib/local';

  let {
    decision,
    inviter,
    onjoined,
    onpeek,
  }: { decision: Decision; inviter: string | null; onjoined: (m: Me) => void; onpeek: () => void } = $props();

  // three short screens: who asked and what for, what you'll do, your name
  const STEPS = ['welcome', 'how', 'name'] as const;
  let step = $state(0);
  let back = $state(false);

  let name = $state(lastName());
  let error = $state('');
  let busy = $state(false);

  const still = typeof matchMedia === 'function' && matchMedia('(prefers-reduced-motion: reduce)').matches;
  const slide = $derived({ x: still ? 0 : back ? -40 : 40, duration: still ? 0 : 220 });

  const points = $derived(decision.features.map((f) => f.label));
  const others = $derived(decision.people.map((p) => p.name).filter((n) => n !== inviter));

  function goTo(n: number) {
    back = n < step;
    error = '';
    step = n;
    window.scrollTo(0, 0);
  }

  // "a, b and c"
  function list(items: string[]): string {
    return items.length < 2 ? items.join('') : `${items.slice(0, -1).join(', ')} and ${items.at(-1)}`;
  }

  async function join(e: SubmitEvent) {
    e.preventDefault();
    if (!name.trim()) return (error = 'Enter your name.');
    busy = true;
    error = '';
    try {
      onjoined(await api.join(decision.id, name.trim()));
    } catch (err) {
      error = (err as Error).message;
    } finally {
      busy = false;
    }
  }
</script>

<main class="mx-auto flex min-h-dvh max-w-md flex-col px-5 pt-5 pb-8">
  <nav class="flex h-10 items-center justify-between">
    {#if step > 0}
      <button class="-ml-2 px-2 text-lg text-slate-500" onclick={() => goTo(step - 1)}>‹ Back</button>
    {:else}
      <span></span>
    {/if}
    <div class="flex gap-2" aria-label="Step {step + 1} of {STEPS.length}">
      {#each STEPS as s, i (s)}
        <span class="h-2 rounded-full transition-all {i === step ? 'w-6 bg-brand-600' : i < step ? 'w-2 bg-brand-500' : 'w-2 bg-slate-300 dark:bg-slate-700'}"></span>
      {/each}
    </div>
  </nav>

  {#key step}
    <div class="flex flex-1 flex-col justify-center gap-8 py-6" in:fly={slide}>
      {#if STEPS[step] === 'welcome'}
        <div class="flex flex-col items-center gap-4 text-center">
          {#if inviter}
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-brand-600 text-4xl font-bold text-white" aria-hidden="true">
              {inviter[0].toUpperCase()}
            </div>
            <p class="text-xl text-slate-600 dark:text-slate-300"><b class="text-slate-900 dark:text-white">{inviter}</b> invited you to help decide</p>
          {:else}
            <img src="/icons/icon-192.png" alt="" class="h-20 w-20 rounded-2xl" />
            <p class="text-xl text-slate-600 dark:text-slate-300">You're invited to help decide</p>
          {/if}
          <h1 class="text-3xl leading-tight font-bold break-words">{decision.title}</h1>
          <p class="text-lg text-slate-600 dark:text-slate-300">
            It takes a few minutes. No account, nothing to install.
          </p>
        </div>

        {#if others.length}
          <p class="text-center text-base text-slate-500">
            Also in: {list(others)}.
          </p>
        {/if}

        <div class="flex flex-col gap-3">
          <button class="btn-primary py-5 text-lg" onclick={() => goTo(1)}>Let's start</button>
          {#if decision.people.length}
            <button class="py-2 text-base text-slate-500 underline" onclick={onpeek}>Just look at the results so far</button>
          {/if}
        </div>
      {:else if STEPS[step] === 'how'}
        <h1 class="text-3xl font-bold">Here's what you'll do</h1>

        <ol class="flex flex-col gap-5">
          <li class="flex gap-4">
            <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-brand-100 text-xl font-bold text-brand-700 dark:bg-brand-700/40 dark:text-brand-100">1</span>
            <div>
              <h2 class="text-xl font-semibold">Say what matters to you</h2>
              <p class="mt-1 text-lg text-slate-600 dark:text-slate-300">
                Tap how much you care about {points.length === 1 ? 'the point' : 'each point'}: {list(points)}.
              </p>
            </div>
          </li>
          <li class="flex gap-4">
            <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-brand-100 text-xl font-bold text-brand-700 dark:bg-brand-700/40 dark:text-brand-100">2</span>
            <div>
              <h2 class="text-xl font-semibold">Rate the options</h2>
              <p class="mt-1 text-lg text-slate-600 dark:text-slate-300">
                {#if decision.options.length}
                  Score {decision.options.length === 1 ? 'the option' : `all ${decision.options.length} options`} from 0 (bad) to 5 (great). You can add your own too.
                {:else}
                  Add the options you're choosing between, then score each from 0 (bad) to 5 (great).
                {/if}
              </p>
            </div>
          </li>
          <li class="flex gap-4">
            <span class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-brand-100 text-xl font-bold text-brand-700 dark:bg-brand-700/40 dark:text-brand-100">3</span>
            <div>
              <h2 class="text-xl font-semibold">See the group's pick</h2>
              <p class="mt-1 text-lg text-slate-600 dark:text-slate-300">Everyone's answers become one ranking. Best match on top.</p>
            </div>
          </li>
        </ol>

        <p class="rounded-2xl bg-brand-50 px-4 py-3 text-base text-brand-700 dark:bg-brand-700/20 dark:text-brand-100">
          Everything saves as you tap. You can change your answers any time.
        </p>

        <button class="btn-primary py-5 text-lg" onclick={() => goTo(2)}>Got it</button>
      {:else}
        <form class="flex flex-col gap-5" onsubmit={join}>
          <label class="flex flex-col gap-3">
            <span class="text-3xl font-bold">What's your name?</span>
            <span class="text-lg text-slate-600 dark:text-slate-300">
              So {inviter ?? 'the others'}{inviter && others.length ? ' and the others' : ''} can see who's who.
            </span>
            <!-- svelte-ignore a11y_autofocus -->
            <input bind:value={name} placeholder="e.g. Alice" maxlength="40" autocomplete="given-name" class="py-4 text-xl" autofocus />
          </label>
          {#if error}<p class="text-base text-red-600">{error}</p>{/if}
          <button class="btn-primary py-5 text-lg" disabled={busy}>{busy ? 'Joining…' : 'Join and start'}</button>
        </form>
      {/if}
    </div>
  {/key}
</main>
