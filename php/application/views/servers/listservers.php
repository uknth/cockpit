<!-- Content Header (Page header) -->
<section class="content-header">
    <h1>
        Widgets
        <small>Preview page</small>
    </h1>
    <ol class="breadcrumb">
        <li><a href="#"><i class="fa fa-dashboard"></i> Home</a></li>
        <li class="active">Widgets</li>
    </ol>
</section>

<!-- Main content -->
<section class="content">
    <h4 class="page-header">
        List of servers
        <small>Small boxes are used for viewing statistics. To create a small box use the class <code>.small-box</code> and mix & match using the <code>bg-*</code> classes.</small>
    </h4>
    <?php foreach($servers as $k => $v){
        if($k % 4 == 0 ){
    ?>
    <div class="row">
        <?php } ?>
        <div class="col-lg-3 col-xs-6">
            <!-- small box -->
            <div class="small-box bg-aqua">
                <div class="inner">
                    <h3>
                        <?php echo $v['name']; ?>
                    </h3>
                    <p>
                        <?php echo $v['ip'].' ('.$v['status'].')';?>
                    </p>
                </div>
                <div class="icon">
                    <i class="ion ion-bag"></i>
                </div>
                <a href="<?php ?>" class="small-box-footer">
                    click to edit <i class="fa fa-edit"></i>
                </a>
            </div>
        </div><!-- ./col -->
    <?php 
        if($k % 4 == 3 ){
    ?>
    </div><!-- /.row -->
    <?php
        }
    }?>
</section><!-- /.content -->