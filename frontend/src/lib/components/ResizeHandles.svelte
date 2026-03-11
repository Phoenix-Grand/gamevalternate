<script lang="ts">
  const THICKNESS = 6;

  function startResize(direction: string) {
    (window as any).chrome?.webview?.postMessage(`resize:${direction}`);
  }
</script>

<!-- Edge handles -->
<div class="handle n"  style="top:0;left:{THICKNESS}px;right:{THICKNESS}px;height:{THICKNESS}px;cursor:n-resize"  onmousedown={() => startResize('n-resize')}></div>
<div class="handle s"  style="bottom:0;left:{THICKNESS}px;right:{THICKNESS}px;height:{THICKNESS}px;cursor:s-resize"  onmousedown={() => startResize('s-resize')}></div>
<div class="handle w"  style="left:0;top:{THICKNESS}px;bottom:{THICKNESS}px;width:{THICKNESS}px;cursor:w-resize"  onmousedown={() => startResize('w-resize')}></div>
<div class="handle e"  style="right:0;top:{THICKNESS}px;bottom:{THICKNESS}px;width:{THICKNESS}px;cursor:e-resize"  onmousedown={() => startResize('e-resize')}></div>

<!-- Corner handles -->
<div class="handle nw" style="top:0;left:0;width:{THICKNESS}px;height:{THICKNESS}px;cursor:nw-resize" onmousedown={() => startResize('nw-resize')}></div>
<div class="handle ne" style="top:0;right:0;width:{THICKNESS}px;height:{THICKNESS}px;cursor:ne-resize" onmousedown={() => startResize('ne-resize')}></div>
<div class="handle sw" style="bottom:0;left:0;width:{THICKNESS}px;height:{THICKNESS}px;cursor:sw-resize" onmousedown={() => startResize('sw-resize')}></div>
<div class="handle se" style="bottom:0;right:0;width:{THICKNESS}px;height:{THICKNESS}px;cursor:se-resize" onmousedown={() => startResize('se-resize')}></div>

<style>
  .handle {
    position: fixed;
    z-index: 9999;
  }
</style>
