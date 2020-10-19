<!-- 预定义的模态框 -->
<div class="modal fade" id="MyModal">
    <div class="modal-dialog {{.}}" id="ModalSize"><!-- 中号模态框 -->
    <!-- div class="modal-dialog modal-lg" --><!-- 大号模态框 -->
    <!-- div class="modal-dialog modal-sm" --><!-- 小号模态框 -->
        <div class="modal-content">
   
            <!-- 模态框头部 -->
            <div class="modal-header">
              <h4 id="MyModal-Title" class="modal-title">模态框头部</h4>
              <button type="button" class="close" data-dismiss="modal">&times;</button>
            </div>
   
            <!-- 模态框主体 -->
            <div id="MyModal-Text" class="modal-body">
              模态框内容..
            </div>
            <div class="modal-body" id='Modal_Echarts' style='height: 500px;width:768px;display:none;'></div>
            <!-- 模态框底部 -->
            <div class="modal-footer">
				      <button type="button" class="btn btn-primary" id="modal-btn" data-dismiss="modal"></button>
            </div>
        </div>
    </div>
</div>