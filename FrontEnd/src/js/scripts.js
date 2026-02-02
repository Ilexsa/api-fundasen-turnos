    // ==========================================
    // CONFIGURACIÓN
    // ==========================================
    const VIDEO_VOL_NORMAL = 1.0;   // volumen habitual
    const VIDEO_VOL_DUCK = 0.15;    // volumen mientras narra (0.1–0.3 recomendado)
    const BASE_HOST = "10.0.0.16:9090"; 
    const WS_URL = (location.protocol === 'https:' ? 'wss' : 'ws') + '://' + BASE_HOST + '/ws';  
    let MI_UBICACION = localStorage.getItem("MI_UBICACION") || "";
    const API_HISTORY_BASE =
    (location.protocol === 'https:' ? 'https' : 'http') + '://' + BASE_HOST + '/api/turnos/estado/';

const selUb = document.getElementById("selUbicacion");
if (selUb) {
  selUb.value = MI_UBICACION;

  selUb.addEventListener("change", () => {
    applyUbicacion(selUb.value);
  });
}

    // =======================
    // UTILIDADES
    // =======================

function getHistoryUrl() {
  const ub = (MI_UBICACION ?? "").trim();
  return ub ? `${API_HISTORY_BASE}${encodeURIComponent(ub)}` : API_HISTORY_BASE;
}



    function applyUbicacion(newUb) {
  MI_UBICACION = (newUb ?? "").trim();
  localStorage.setItem("MI_UBICACION", MI_UBICACION);

  // limpiar cola para que no se cuelen turnos de otra ubicación
  state.queue = [];

  // si el current no cumple, lo borramos
  if (state.current && !matchUbicacion(state.current)) {
    state.current = null;
  }

  // filtrar historial existente
  state.history = (state.history || []).filter(matchUbicacion);

  render();   // refresco inmediato

  // traer verdad de BD (ya filtrada en front)
  init();
}
    function toTitleCase(str) {
      if (!str) return "";
      return String(str).toLowerCase().split(' ').map(word => 
        word.charAt(0).toUpperCase() + word.slice(1)
      ).join(' ');
    }



    // Normaliza el payload entrante (WS/API) a una estructura interna única
    // ✅ Internamente todo usa: id_consulta
    function normalizeTurno(t) {
      if (!t) return null;

      const id = t.id_consulta ?? t.consulta_id ?? t.consultaId ?? t.consultaID ?? null;
      if (id === null || id === undefined || String(id).trim() === "") return null;

      return {
        id_consulta: String(id),
        paciente: t.paciente ?? "",
        medico: t.medico ?? "",
        consultorio: t.consultorio ?? "",
        especialidad: t.especialidad ?? "",
        localidad: t.localidad ?? "",
        ubicacion: t.ubicacion ?? ""
      };
    }

    function normTxt(s){
  return String(s ?? "").trim().toLowerCase().replace(/\s+/g, " ");
}

function matchUbicacion(turno){
  if (!MI_UBICACION || normTxt(MI_UBICACION) === "general") return true;

  const uTurno = normTxt(turno?.ubicacion);
  const uPantalla = normTxt(MI_UBICACION);

  // Si el backend manda "General", se ve en todas
  if (uTurno === "general") return true;

  return uTurno === uPantalla;
}

    // =======================
    // RELOJ
    // =======================
    function tick(){
      const now = new Date();
      document.getElementById('time').textContent = now.toLocaleTimeString('es-EC', {hour:'2-digit', minute:'2-digit'});
      document.getElementById('date').textContent = now.toLocaleDateString('es-EC', {weekday:'long', year:'numeric', month:'long', day:'numeric'});
    }
    tick(); setInterval(tick, 1000);

    // =======================
    // VIDEOS
    // =======================
const videoEl = document.getElementById('promoVideo');
let promoIndex = 0;
let srcIndex = 0;

const promos = [
    {
    badge: "FUNDASEN INFORMA",
    title: "REFLEXIÓN",
    sub: "",
    srcs: [
        "./src/media/videos/VIDEO NUEVO TURNERO 2026.mp4"
    ]
    },
/*        {
    badge: "TWENTY ONE PILOTS COVER",
    title: "SEVEN NATION ARMY",
    sub: "COVER",
    srcs: [
        "./src/media/videos/Twenty One Pilots Cover The White Stripes _Seven Nation Army_ _ Rock Hall 2025 Induction.mp4"
    ]
    },
    {
    badge: "CAN'T STOP",
    title: "RED HOT CHILLI PEAPERS",
    sub: "CUALQUIER TEXTO DE SUB TITULO",
    srcs: [
        "./src/media/videos/Red Hot Chili Peppers - Can't Stop [Official Music Video].mp4",
    ]
    },
    {
    badge: "Freak On a Leash",
    title: "Korn",
    sub: "Pruebas",
    srcs: [
        "./src/media/videos/Korn - Freak On a Leash (Official HD Video).mp4"
    ]
    }*/
];

function setPromoUI(p){
  document.getElementById('videoBadge').textContent = p.badge || "Info";
  document.getElementById('videoTitle').textContent = p.title || "Video";
  document.getElementById('videoSub').textContent = p.sub || "";
}

function playNextVideo(){
  if (!promos.length) return;

  const p = promos[promoIndex];
  setPromoUI(p);

  const list = p.srcs || [];
  if (!list.length) return;

  videoEl.src = list[srcIndex];
  videoEl.play().catch(()=>{});

  // avanzar índices
  srcIndex++;
  if (srcIndex >= list.length) {
    srcIndex = 0;
    promoIndex = (promoIndex + 1) % promos.length;
  }
}
function playPrevVideo(){
  if (!promos.length) return;

  // retroceder dentro del promo actual o al promo anterior
  srcIndex--;
  if (srcIndex < 0) {
    promoIndex = (promoIndex - 1 + promos.length) % promos.length;
    const prevPromo = promos[promoIndex];
    const list = prevPromo.srcs || [];
    srcIndex = Math.max(0, list.length - 1);
  }

  const p = promos[promoIndex];
  setPromoUI(p);

  const list = p.srcs || [];
  if (!list.length) return;

  videoEl.src = list[srcIndex];
  videoEl.play().catch(()=>{});
}

videoEl.addEventListener('ended', playNextVideo);
videoEl.addEventListener('error', playNextVideo);

// arranque
playNextVideo();
window.promoNext = playNextVideo;
window.promoPrev = playPrevVideo;
    // =======================
    // ESTADO GLOBAL
    // =======================
    const state = {
      current: null,      // El turno grande actual
      history: [],        // La lista de la derecha
      queue: [],          // COLA DE ESPERA (Buffer)
      isBusy: false       // SEMÁFORO: True si está hablando o procesando
    };

    // =======================
    // RENDERIZADO
    // =======================
    function render() {
      // A. Render Turno Grande
      if (state.current) {
          document.getElementById('currentPatient').textContent = state.current.paciente;
          document.getElementById('currentRoom').textContent = state.current.consultorio;
          document.getElementById('currentExtra').textContent = state.current.especialidad;
          document.getElementById('currentUbi').textContent = state.current.ubicacion;
      }

      // B. Render Historial (Lista Lateral)
      const list = document.getElementById('queueList');
      list.innerHTML = "";
      
      // Filtramos para que el turno ACTUAL no salga repetido en el historial visual
      const historialLimpio = state.history.filter(t => {
          if (!state.current) return true;
          return t.id_consulta !== state.current.id_consulta;
      });

      historialLimpio.slice(0, 10).forEach(p => {
        const item = document.createElement('div');
        item.className = 'queue-item';
        item.innerHTML = `
          <div>
            <p class="name">${p.paciente}</p>
            <p class="meta">${p.especialidad}</p>
          </div>
          <div class="room-tag">${p.consultorio}</div>
          <div class="ubi-tag">${p.ubicacion}</div>
        `;
        list.appendChild(item);
      });
      document.getElementById('queueCount').textContent = historialLimpio.length;
    }

    // =======================
    // GESTOR DE COLA (QUEUE MANAGER)
    // =======================
    
    // Función principal que intenta procesar el siguiente elemento
    async function procesarCola() {
        // 1. REGLA DE ORO: Si está ocupado o no hay nada, no hacer nada.
        if (state.isBusy) return;
        if (state.queue.length === 0) return;

        // 2. BLOQUEAR EL SISTEMA
        state.isBusy = true;

        // 3. OBTENER DATOS
        const nuevoTurno = state.queue.shift(); // Sacar el primero de la fila (FIFO)

        // 4. ACTUALIZACIÓN VISUAL INMEDIATA (Fase 1)
        // Movemos el actual al historial localmente para que se vea fluido
        if (state.current) {
            state.history.unshift(state.current);
        }
        state.current = nuevoTurno;
        render(); // Pinta la pantalla

        // 5. LLAMADA API (GET) - Como solicitaste
        // Hacemos el fetch AHORA para asegurar que el historial esté sincronizado con BD
        try {
            console.log("Fetching historial actualizado..."); 
            const res = await fetch(getHistoryUrl());
            if (res.ok) {
                const data = await res.json();
                if (Array.isArray(data)) {
                    //state.history = data.map(normalizeTurno).filter(Boolean); // verdad de BD, normalizada
                    state.history = data.map(normalizeTurno).filter(Boolean).filter(matchUbicacion);
                    render(); // Repintamos la lista lateral
                }
            }
        } catch (e) {
            console.error("Error fetching history", e);
        }

        // 6. INICIAR AUDIO (Esto bloqueará el siguiente turno hasta que termine)
        await hablar(nuevoTurno);

        // 7. DESBLOQUEAR Y REVISAR SI HAY MÁS
        state.isBusy = false;
        
        // Pequeño delay para que no sea atropellado
        setTimeout(() => {
            procesarCola(); // Recursividad: Llamar al siguiente si existe
        }, 1000);
    }

    // =======================
    // AUDIO (PROMISIFIED)
    // =======================
    let audioEnabled = false;
    let voices = [];
    if ('speechSynthesis' in window) {
        window.speechSynthesis.onvoiceschanged = () => voices = window.speechSynthesis.getVoices();
    }


    function hablar(turno) {
  return new Promise((resolve) => {
    // UI visual de "Llamando"
    const badge = document.getElementById('callBadge');
    const card = document.getElementById('currentCard');
    badge.classList.add('show');
    card.classList.add('calling');
    document.getElementById('callText').textContent = "LLAMANDO...";

    // VIDEO: ducking con fade
    const videoEl = document.getElementById('promoVideo');
    const prevVol = (videoEl && typeof videoEl.volume === "number") ? videoEl.volume : 1.0;
    const DUCK_VOL = 0.05;     // volumen mientras narra
    const FADE_MS = 300;       // duración del fade
    const STEP_MS = 25;        // suavidad

    let fadeTimer = null;
    let finished = false;

    function fadeTo(targetVol) {
      if (!videoEl || videoEl.muted) return;
      if (fadeTimer) clearInterval(fadeTimer);

      const startVol = videoEl.volume;
      const steps = Math.max(1, Math.round(FADE_MS / STEP_MS));
      const delta = (targetVol - startVol) / steps;
      let i = 0;

      fadeTimer = setInterval(() => {
        i++;
        let next = videoEl.volume + delta;

        // clamp 0..1
        if (next < 0) next = 0;
        if (next > 1) next = 1;

        videoEl.volume = next;

        if (i >= steps) {
          videoEl.volume = targetVol;
          clearInterval(fadeTimer);
          fadeTimer = null;
        }
      }, STEP_MS);
    }

    // Bajar volumen ANTES de hablar
    showCallOverlay(turno);
    fadeTo(DUCK_VOL);

    const cleanup = () => {
      if (finished) return; // evita doble cleanup por onend/onerror/cancel
      finished = true;


      hideCallOverlay();
      // Restaurar volumen con fade
      fadeTo(prevVol);

      badge.classList.remove('show');
      card.classList.remove('calling');
      resolve();
    };

    // Si no hay audio habilitado, esperamos 4 seg y continuamos
    if (!audioEnabled || !('speechSynthesis' in window)) {
      setTimeout(cleanup, 4000);
      return;
    }

    // Configurar voz
    const txt = `Paciente ${toTitleCase(turno.paciente)}, pasar al ${toTitleCase(turno.consultorio)}, ${toTitleCase(turno.especialidad)}`;
    const u = new SpeechSynthesisUtterance(txt);
    u.lang = "es-EC";
    u.volume = 1;
    u.rate = 0.9;

    const v = voices.find(vo => vo.lang.includes('es') && vo.name.includes('Google'))
            || voices.find(vo => vo.lang.includes('es'));
    if (v) u.voice = v;

    // Eventos para liberar la promesa
    u.onend = cleanup;
    u.onerror = cleanup;

    // Cancelar cualquier cosa previa y hablar
    window.speechSynthesis.cancel();
    window.speechSynthesis.speak(u);
  });
}

    // Botón activar audio
document.getElementById('btnEnableAudio').addEventListener('click', async () => {
    // 1) Tomar ubicación desde el combo
  const selUb = document.getElementById("selUbicacion");
  MI_UBICACION = (selUb?.value ?? "").trim();  // "" => general
  localStorage.setItem("MI_UBICACION", MI_UBICACION);
  
  audioEnabled = true;
  document.getElementById('audioGate').style.display = 'none';

  // ✅ Habilitar audio del video
  const videoEl = document.getElementById('promoVideo');
  videoEl.muted = false;
  videoEl.volume = 0.5;

  try { 
    await videoEl.play(); // por si el navegador pausó el autoplay
  } catch (e) {
    console.log("No se pudo reproducir el video con audio:", e);
  }
  applyUbicacion(selUb?.value || "");
  document.getElementById('statusBadge').textContent = "En línea • " + (MI_UBICACION ? MI_UBICACION : "GENERAL");
});


    // =======================
    // CARGA INICIAL
    // =======================
    async function init() {
        try {
            const res = await fetch(getHistoryUrl());
            const data = await res.json();
            if (Array.isArray(data) && data.length > 0) {
                //const norm = data.map(normalizeTurno).filter(Boolean);
                const norm = data.map(normalizeTurno).filter(Boolean).filter(matchUbicacion);
                state.current = norm[0] || null;
                state.history = norm.slice(1);
                render();
            }
        } catch (e) { console.error("Error init:", e); }
    }


    // =======================
    // WEBSOCKET
    // =======================
    let socket;
    function connectWS() {
        console.log("Conectando WS...", WS_URL);
        socket = new WebSocket(WS_URL);
        
        socket.onopen = () => {
            document.getElementById('statusBadge').textContent = "En línea • " + (MI_UBICACION ? MI_UBICACION : "GENERAL");
            document.getElementById('statusBadge').style.color = "green";
            init(); // Cargar estado inicial
        };

        socket.onmessage = (event) => {
            try {
                const turno = normalizeTurno(JSON.parse(event.data));
                if (!turno || !turno.id_consulta) return;
                //if (MI_UBICACION && turno.ubicacion !== MI_UBICACION && turno.ubicacion !== 'General') return;
                if (!matchUbicacion(turno)) return;

                // FILTRO DE REPETIDOS (Doble clic):
                // Verificamos si ya está en cola o si es el actual
                const yaEnCola = state.queue.find(t => t.id_consulta === turno.id_consulta);
                const esActual = state.current && state.current.id_consulta === turno.id_consulta;
                
                if (!yaEnCola && !esActual) {
                    console.log("Nuevo turno encolado:", turno.paciente);
                    state.queue.push(turno); // 1. Solo pushear a la cola
                    procesarCola();          // 2. Intentar arrancar el procesador
                } else {
                    console.log("Turno ignorado (duplicado/actual)");
                }

            } catch (e) { console.error(e); }
        };

        socket.onclose = () => {
            document.getElementById('statusBadge').textContent = "Desconectado";
            document.getElementById('statusBadge').style.color = "red";
            setTimeout(connectWS, 3000);
        };
    }
    function showCallOverlay(turno){
  const ov = document.getElementById('callOverlay');
  if (!ov) return;

  document.getElementById('callPatient').textContent = turno.paciente;
  document.getElementById('callRoom').textContent = turno.consultorio;
  document.getElementById('callSpec').textContent = turno.especialidad;
  document.getElementById('callUbi').textContent = `${turno.ubicacion}`;
  document.getElementById('callUb').textContent = `Ubicación: ${toTitleCase(turno.ubicacion)}`;

  ov.classList.add('show');


  const ding = document.getElementById('callDing');
  if (ding && !ding.muted) {
    try { ding.currentTime = 0; ding.play().catch(()=>{}); } catch {}
  }
}

function hideCallOverlay(){
  const ov = document.getElementById('callOverlay');
  if (!ov) return;
  ov.classList.remove('show');
}

(() => {
  const video  = document.getElementById("promoVideo");
  const stage  = document.getElementById("videoStage");

  const btnPlay = document.getElementById("btnPlay");
  const btnMute = document.getElementById("btnMute");
  const btnFs   = document.getElementById("btnFs");

  const seek    = document.getElementById("seek");
  const vol     = document.getElementById("vol");

  const timeNow   = document.getElementById("timeNow");
  const timeTotal = document.getElementById("timeTotal");

  const btnPrevVid = document.getElementById("btnPrevVid");
  const btnNextVid = document.getElementById("btnNextVid");

  // Helpers
  const fmt = (sec) => {
    if (!isFinite(sec)) return "0:00";
    sec = Math.max(0, Math.floor(sec));
    const m = Math.floor(sec / 60);
    const s = String(sec % 60).padStart(2, "0");
    return `${m}:${s}`;
  };

  const setPlayIcon = () => {
    if (btnPlay) btnPlay.textContent = video.paused ? "▶" : "⏸";
  };

  const setMuteIcon = () => {
    if (btnMute) btnMute.textContent = (video.muted || video.volume === 0) ? "🔇" : "🔊";
  };

  const resetUI = () => {
    if (seek) seek.value = "0";
    if (timeNow) timeNow.textContent = "0:00";
    if (timeTotal) timeTotal.textContent = "0:00";
  };

  // ---- Play/Pausa
  btnPlay?.addEventListener("click", () => {
    if (video.paused) video.play().catch(() => {});
    else video.pause();
  });

  video.addEventListener("play", setPlayIcon);
  video.addEventListener("pause", setPlayIcon);
  setPlayIcon();

  // ---- Mute/Volumen
  btnMute?.addEventListener("click", () => {
    video.muted = !video.muted;
    setMuteIcon();
  });

  vol?.addEventListener("input", () => {
    video.volume = Number(vol.value);
    if (video.volume > 0) video.muted = false;
    setMuteIcon();
  });

  // sincroniza UI inicial
  if (vol) vol.value = video.volume;
  setMuteIcon();

  // ---- Duración y barra
  video.addEventListener("loadedmetadata", () => {
    if (timeTotal) timeTotal.textContent = fmt(video.duration);
  });

  video.addEventListener("timeupdate", () => {
    if (timeNow) timeNow.textContent = fmt(video.currentTime);
    if (seek && video.duration) {
      const pct = (video.currentTime / video.duration) * 100;
      seek.value = String(pct);
    }
  });

  // ---- Seek
  const seekTo = () => {
    if (!seek || !video.duration) return;
    const pct = Number(seek.value) / 100;
    video.currentTime = pct * video.duration;
  };

  let seeking = false;
  seek?.addEventListener("pointerdown", () => seeking = true);
  seek?.addEventListener("pointerup",   () => { seeking = false; seekTo(); });
  seek?.addEventListener("input", () => { if (seeking) seekTo(); });

  // ---- Fullscreen
  btnFs?.addEventListener("click", async () => {
    try {
      if (!document.fullscreenElement) await stage.requestFullscreen();
      else await document.exitFullscreen();
    } catch (_) {}
  });

  // ---- Auto-hide controles
  let hideT = null;
  const showControls = () => {
    stage.classList.remove("controls-hidden");
    if (hideT) clearTimeout(hideT);
    hideT = setTimeout(() => stage.classList.add("controls-hidden"), 1800);
  };

  stage.addEventListener("mousemove", showControls);
  stage.addEventListener("mouseenter", showControls);
  stage.addEventListener("mouseleave", () => stage.classList.remove("controls-hidden"));
  stage.addEventListener("touchstart", showControls, { passive: true });

  // ============================================================
  //          PREV / NEXT VIDEO (USA TU COLA YA CARGADA)
  //  Requisitos: window.playlist = [url1, url2, ...]
  //              window.videoIndex = índice actual (opcional)
  // ============================================================

  const getPlaylist = () => Array.isArray(window.playlist) ? window.playlist : [];
  const getIndex = () => Number.isFinite(Number(window.videoIndex)) ? Number(window.videoIndex) : 0;
  const setIndex = (i) => { window.videoIndex = i; };

  const changeVideoTo = (newIndex) => {
    const list = getPlaylist();
    if (!list.length) return;

    // conservar audio elegido por usuario
    const keepMuted = video.muted;
    const keepVol = video.volume;

    const len = list.length;
    const idx = (newIndex + len) % len;     // loop entre videos
    setIndex(idx);

    resetUI();

    video.src = list[idx];
    video.load();

    video.muted = keepMuted;
    video.volume = keepVol;
    if (vol) vol.value = keepVol;
    setMuteIcon();

    video.play().catch(() => {});
  };

  const nextVideo = () => changeVideoTo(getIndex() + 1);
  const prevVideo = () => changeVideoTo(getIndex() - 1);

btnNextVid?.addEventListener("click", () => { resetUI(); window.promoNext?.(); });
btnPrevVid?.addEventListener("click", () => { resetUI(); window.promoPrev?.(); });

  // Si quieres que al terminar pase al siguiente de la cola:
  video.addEventListener("ended", () => {
    // Si tu loop ya lo maneja otro código, comenta esta línea:
    nextVideo();
  });

  // ---- Inicial
  showControls();

  // Si quieres que al cargar la página use el índice actual de tu cola:
  // (solo si playlist ya existe antes de ejecutar este script)
  const list = getPlaylist();
  if (list.length) {
    // si ya tenías src puesto y no quieres recargar, comenta esto:
    changeVideoTo(getIndex());
  }
})();

(() => {
  const secret = ["ArrowUp","ArrowUp","ArrowDown","ArrowRight","ArrowRight"];
  let idx = 0;
  let timer = null;

  const img = document.getElementById("toastyEgg");
  const audio = document.getElementById("toastyAudio");

  function triggerToasty(){
    // GIF
    img.classList.remove("show"); // reinicia animación
    void img.offsetWidth;         // reflow para que vuelva a animar
    img.classList.add("show");

    // AUDIO (algunos navegadores solo dejan si hubo interacción antes)
    try{
      audio.currentTime = 0;
      audio.play();
    }catch(e){}
  }

  window.addEventListener("keydown", (e) => {
    const key = e.key;

    // opcional: evita disparos si estás escribiendo en inputs
    const tag = (document.activeElement?.tagName || "").toLowerCase();
    if (tag === "input" || tag === "textarea") return;

    // reset por timeout (si te demoras, se reinicia)
    clearTimeout(timer);
    timer = setTimeout(() => (idx = 0), 1200);

    if (key === secret[idx]) {
      idx++;
      if (idx === secret.length) {
        idx = 0;
        triggerToasty();
      }
    } else {
      idx = 0;
    }
  });
})();

    connectWS();